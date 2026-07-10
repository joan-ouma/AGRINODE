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

type SensorData struct {
	Temperature  float32 `json:"temperature"`
	Humidity     float32 `json:"humidity"`
	SoilMoisture int     `json:"soilMoisture"`
}

var kafkaProducer *kafka.Producer
var kafkaTopic = "raw-telemetry-stream"

func main() {
	var err error

	// 1. Init Kafka Producer
	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "agrinode-kafka:9092"})
	if err != nil {
		log.Fatal("failed to create kafka producer: ", err)
	}
	defer kafkaProducer.Close()
	fmt.Println("kafka producer initialized.")

	// Handle delivery reports in the background
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

	// 2. Init MQTT Subscriber
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://agrinode-broker:1883")
	opts.SetClientID("agrinode-ingest-kafka")

	// THESE MATCH MOSQUITTO SETUP
	opts.SetUsername("agrinode_device")
	opts.SetPassword("farm_secret")

	opts.SetDefaultPublishHandler(messageHandler)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("mqtt connection failed: ", token.Error())
	}
	fmt.Println("connected to mqtt broker.")

	topic := "agrinode/telemetry"
	if token := mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal("mqtt subscribe failed: ", token.Error())
	}
	fmt.Printf("listening on %s -> forwarding to kafka topic: %s\n", topic, kafkaTopic)

	// Wait for Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nshutting down ingestion service...")
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	var data SensorData
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("json parse error: %v\n", err)
		return
	}

	// Simple validation to ensure clean data
	if data.Humidity < 0 || data.Humidity > 100 {
		return
	}

	// Push raw payload to Kafka
	err := kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          msg.Payload(),
	}, nil)

	if err != nil {
		log.Printf("failed to enqueue message: %v\n", err)
		return
	}

	fmt.Printf("pushed to kafka -> temp: %.1f, humidity: %.1f, moisture: %d\n", data.Temperature, data.Humidity, data.SoilMoisture)
}
