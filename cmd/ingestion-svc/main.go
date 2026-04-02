package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SensorData struct perfectly matches the JSON payload from the ESP8266
type SensorData struct {
	Temperature  float32 `json:"temperature"`
	Humidity     float32 `json:"humidity"`
	SoilMoisture int     `json:"soilMoisture"`
}

func handleSensorData(w http.ResponseWriter, r *http.Request) {
	// 1. Security Check: Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Catch and Decode the JSON Payload
	var data SensorData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 3. Log the successful ingestion to the terminal
	fmt.Printf("🌱 [New Data Received] Temp: %.2f°C | Hum: %.2f%% | Soil: %d%%\n",
		data.Temperature, data.Humidity, data.SoilMoisture)

	// 4. Send a 200 OK response back to the ESP8266 so it knows the data arrived safely
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received successfully!"))
}

func main() {
	fmt.Println("==========================================")
	fmt.Println("  Agri-Node Ingestion Service Starting    ")
	fmt.Println("==========================================")

	// Route incoming traffic on /api/sensor-data to our handler function
	http.HandleFunc("/api/sensor-data", handleSensorData)

	// Start the server on port 8080
	port := ":8080"
	fmt.Printf("Server listening on port %s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
