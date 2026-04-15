package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
)

type SensorData struct {
	Temperature  float32 `json:"temperature"`
	Humidity     float32 `json:"humidity"`
	SoilMoisture int     `json:"soilMoisture"`
}

func main() {
	// 1. Connect to PostgreSQL (Updated with your exact compose credentials)
	connStr := "postgres://agrinode_admin:supersecretpassword@localhost:5432/agrinode?sslmode=disable"
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open db connection: ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("failed to ping db: ", err)
	}
	fmt.Println("connected to postgresql.")

	// 2. Connect to Kafka Consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "agrinode-storage-group",
		"auto.offset.reset": "earliest", 
	})
	if err != nil {
		log.Fatal("failed to create kafka consumer: ", err)
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"raw-telemetry-stream"}, nil)
	fmt.Println("listening to kafka topic: raw-telemetry-stream...")

	// 3. The Consume Loop
	run := true
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Background thread to listen for Ctrl+C
	go func() {
		<-sigchan
		run = false
	}()

	for run {
		// Wait for a message
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var data SensorData
			// Parse the JSON
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("json parse error: %v\n", err)
				continue
			}

			// Insert into the database
			query := `INSERT INTO telemetry (temperature, humidity, soil_moisture) VALUES ($1, $2, $3)`
			_, err = db.Exec(query, data.Temperature, data.Humidity, data.SoilMoisture)
			
			if err != nil {
				log.Printf("db insert failed: %v\n", err)
			} else {
				fmt.Printf("saved to db -> temp: %.1f, humidity: %.1f, moisture: %d\n", data.Temperature, data.Humidity, data.SoilMoisture)
			}
		} else {
			log.Printf("kafka read error: %v\n", err)
		}
	}

	fmt.Println("\nshutting down storage service...")
}
