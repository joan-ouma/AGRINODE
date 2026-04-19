package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Telemetry struct (Updated to remove Soil Moisture and add Zone ID)
type SensorData struct {
	ZoneID      int     `json:"zoneId"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

// Security payload from the HMI/Keypad
type AuthRequest struct {
	PIN string `json:"pin"`
}

var kafkaProducer *kafka.Producer
var kafkaTopic = "raw-telemetry-stream"

// --- NEW: Security State Variables ---
var failedAttempts int = 0
var systemLocked bool = false

const validPIN = "1234" // Hardcoded for prototype, usually checked against DB

// --- NEW: Telegram Configuration ---
const telegramToken = "YOUR_BOT_TOKEN_HERE"
const chatID = "YOUR_CHAT_ID_HERE"

func main() {
	var err error

	// 1. Init Kafka Producer
	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatal("failed to create kafka producer: ", err)
	}
	defer kafkaProducer.Close()
	fmt.Println("✅ Kafka producer initialized.")

	go func() {
		for e := range kafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("delivery failed: %v\n", ev.TopicPartition.Error)
				}
			}
		}
	}()

	// 2. Init MQTT Subscriber & Publisher
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("bms-core-ingestion")
	opts.SetUsername("agrinode_device")
	opts.SetPassword("farm_secret")

	// Set the default handler, but we will route based on topic inside it
	opts.SetDefaultPublishHandler(messageHandler)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("mqtt connection failed: ", token.Error())
	}
	fmt.Println("✅ Connected to MQTT broker.")

	// Subscribe to BOTH Telemetry and Security topics
	topics := map[string]byte{
		"bms/telemetry":     0,
		"bms/security/auth": 0,
	}
	if token := mqttClient.SubscribeMultiple(topics, nil); token.Wait() && token.Error() != nil {
		log.Fatal("mqtt subscribe failed: ", token.Error())
	}
	fmt.Println("🎧 Listening for Telemetry and Security events...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down ingestion service...")
}

// The master router for incoming MQTT messages
func messageHandler(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()

	if topic == "bms/telemetry" {
		handleTelemetry(msg.Payload())
	} else if topic == "bms/security/auth" {
		handleSecurityAuth(client, msg.Payload())
	}
}

// --- NEW: Security & 3-Strike Logic ---
func handleSecurityAuth(client mqtt.Client, payload []byte) {
	var req AuthRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Println("Invalid auth payload")
		return
	}

	if systemLocked {
		fmt.Println("🔒 System is locked. Ignoring input.")
		return
	}

	if req.PIN == validPIN {
		// SUCCESS
		failedAttempts = 0
		fmt.Println("✅ Entry Granted. Unlocking Control Room.")

		// Command the DFPlayer to play Track 1 ("Entry Granted")
		client.Publish("bms/control-room/audio", 0, false, "play_track_1")
		// Command the Servo to unlock the door
		client.Publish("bms/control-room/door", 0, false, "unlock")

	} else {
		// FAILURE
		failedAttempts++
		fmt.Printf("❌ Wrong PIN. Attempt %d/3\n", failedAttempts)

		if failedAttempts >= 3 {
			systemLocked = true
			fmt.Println("🚨 SYSTEM LOCKDOWN INITIATED!")

			// Command the buzzer/DFPlayer to sound the continuous alarm (Track 3)
			client.Publish("bms/control-room/audio", 0, false, "play_track_3")

			// Fire the Cloud Alert
			go sendTelegramAlert("🚨 URGENT: 3 consecutive failed login attempts detected at the Control Room. System is currently in Lockdown.")
		} else {
			// Just a warning
			client.Publish("bms/control-room/audio", 0, false, "play_track_2") // "Wrong Attempt"
		}
	}
}

// Existing Telemetry logic, updated to push the new Zone data to Kafka
func handleTelemetry(payload []byte) {
	var data SensorData
	if err := json.Unmarshal(payload, &data); err != nil {
		return
	}

	err := kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)

	if err != nil {
		log.Printf("failed to enqueue telemetry: %v\n", err)
		return
	}
	fmt.Printf("📊 Pushed to Kafka -> Zone: %d, Temp: %.1f\n", data.ZoneID, data.Temperature)
}

// --- NEW: Telegram Alerting Function ---
func sendTelegramAlert(message string) {
	if telegramToken == "YOUR_BOT_TOKEN_HERE" {
		fmt.Println("⚠️ Telegram alert skipped (Token not configured)")
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken)
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Failed to send Telegram alert: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("📲 Telegram Alert Delivered!")
}
