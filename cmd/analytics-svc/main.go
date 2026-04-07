package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// same struct so the shape of our data remains consistent
type SensorData struct {
	Temperature  float32   `json:"temperature"`
	Humidity     float32   `json:"humidity"`
	SoilMoisture int       `json:"soilMoisture"`
	Timestamp    time.Time `json:"timestamp"`
}

var collection *mongo.Collection

func main() {
	// 1. Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("failed to connect to mongo: ", err)
	}
	defer mongoClient.Disconnect(ctx)

	collection = mongoClient.Database("agrinode").Collection("telemetry")
	fmt.Println("Analytics API connected to MongoDB.")

	// 2. Set up the HTTP Routes
	http.HandleFunc("/api/history", getTelemetryHistory)

	// 3. Start the Server
	port := ":8080"
	fmt.Printf("API Server is running and listening on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server crashed: ", err)
	}
}

// getTelemetryHistory pulls the latest 20 readings from the database
func getTelemetryHistory(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers so a React frontend running on a different port can access this API
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Sort by timestamp descending (-1) to get the newest data first
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	findOptions.SetLimit(20)

	// Empty bson.D{} means "match everything"
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []SensorData
	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, "Failed to decode data", http.StatusInternalServerError)
		return
	}

	// If the database is completely empty, return an empty array instead of null
	if results == nil {
		results = []SensorData{}
	}

	// Send the JSON response back to the browser/frontend
	json.NewEncoder(w).Encode(results)
}
