package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Anomaly struct matches the JSON broadcasted by PostgreSQL
type Anomaly struct {
	ID            int     `json:"id"`
	NodeID        int     `json:"node_id"`
	AnomalyType   string  `json:"anomaly_type"`
	CurrentValue  float32 `json:"current_value"`
	MovingAverage float32 `json:"moving_average"`
	DetectedAt    string  `json:"detected_at"`
}

func main() {
	// 1. Connect to the MQTT Broker (To talk to the ESP8266)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("agrinode_command_svc")
	opts.SetUsername("agrinode_device")
	opts.SetPassword("farm_secret")

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("failed to connect to mqtt broker: ", token.Error())
	}
	log.Println("command service connected to mqtt broker...")

	// 2. Connect to Kafka (To listen for alarms)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "command-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal("failed to create kafka consumer: ", err)
	}
	defer consumer.Close()

	// Subscribe to the topic our publisher is sending to
	consumer.Subscribe("anomaly-events", nil)
	log.Println("listening for alarms on kafka topic 'anomaly-events'...")

	// 3. Graceful Shutdown Setup
	run := true
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigchan
		run = false
	}()

	// 4. The Infinite Listening Loop
	for run {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var anomaly Anomaly
			
			// Decode the JSON payload
			if err := json.Unmarshal(msg.Value, &anomaly); err != nil {
				log.Printf("failed to parse anomaly json: %v", err)
				continue
			}

			log.Printf("processing alarm for Node %d: %s", anomaly.NodeID, anomaly.AnomalyType)

			// 5. The Decision Engine: What action do we take?
			var commandPayload string
			
			switch anomaly.AnomalyType {
			case "Critically Dry Soil":
				// Command the water pump to turn on for 60 seconds
				commandPayload = `{"action": "pump_on", "duration_sec": 60}`
			case "High Temperature Spike":
				// Command the cooling fans or shade cloth to activate
				commandPayload = `{"action": "cooling_on", "duration_sec": 300}`
			case "Sudden Humidity Drop":
				// Command the misting system
				commandPayload = `{"action": "misters_on", "duration_sec": 120}`
			default:
				log.Printf("unknown anomaly type, taking no action: %s", anomaly.AnomalyType)
				continue
			}

			// 6. Fire the Command to the specific Edge Node via MQTT
			// Notice how we use the NodeID to dynamically route the message to the right device
			targetTopic := fmt.Sprintf("agrinode/commands/%d", anomaly.NodeID)
			
			token := mqttClient.Publish(targetTopic, 0, false, commandPayload)
			token.Wait()
			
			log.Printf("fired MQTT command to %s -> %s", targetTopic, commandPayload)
		}
	}
}
