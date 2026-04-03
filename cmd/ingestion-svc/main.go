package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorData struct {
	Temperature  float32   `json:"temperature"`
	Humidity     float32   `json:"humidity"`
	SoilMoisture int       `json:"soilMoisture"`
	Timestamp    time.Time `json:"timestamp"`
}

var collection *mongo.Collection

func main() {
	// 1. connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("failed to connect to mongo: ", err)
	}
	defer mongoClient.Disconnect(ctx)

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("mongo ping failed: ", err)
	}
	fmt.Println("connected to mongodb.")

	// set up our database and collection
	collection = mongoClient.Database("agrinode").Collection("telemetry")

	// 2. set up mqtt subscriber
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("agrinode-ingest-svc")
	opts.SetUsername("agrinode_device")
	opts.SetPassword("farm_secret")

	// assign the callback function for when a message arrives
	opts.SetDefaultPublishHandler(messageHandler)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("failed to connect to mqtt: ", token.Error())
	}
	fmt.Println("connected to mqtt broker.")

	// subscribe to the topic
	topic := "agrinode/telemetry"
	if token := mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal("failed to subscribe: ", token.Error())
	}
	fmt.Printf("listening on topic: %s\n", topic)

	// keep the service running until we hit ctrl+c
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nshutting down ingestion service...")
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	var data SensorData

	// parse the incoming json
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Printf("error parsing json: %v\n", err)
		return
	}

	// stamp it with the exact time it arrived
	data.Timestamp = time.Now()

	// basic validation rule: throw out garbage data
	if data.Humidity < 0 || data.Humidity > 100 {
		log.Printf("validation failed: impossible humidity reading %f\n", data.Humidity)
		return
	}

	// save to database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Printf("failed to insert data: %v\n", err)
		return
	}

	fmt.Printf("saved to db -> temp: %.1f, humidity: %.1f, moisture: %d\n", data.Temperature, data.Humidity, data.SoilMoisture)
}
