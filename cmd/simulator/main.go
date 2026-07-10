package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
	Temperature  float32 `json:"temperature"`
	Humidity     float32 `json:"humidity"`
	SoilMoisture int     `json:"soilMoisture"`
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://agrinode-broker:1883")
	opts.SetClientID("agrinode-sim-01")

	// Add these two lines so the simulator can log in:
	opts.SetUsername("agrinode_device")
	opts.SetPassword("farm_secret")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Println("Simulator connected to MQTT broker. Publishing data...")

	for {
		data := SensorData{
			Temperature:  20.0 + rand.Float32()*15.0,
			Humidity:     40.0 + rand.Float32()*20.0,
			SoilMoisture: rand.Intn(100),
		}

		payload, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			continue
		}

		token := client.Publish("agrinode/telemetry", 0, false, payload)
		token.Wait()

		fmt.Printf("Published to agrinode/telemetry: %s\n", string(payload))
		time.Sleep(3 * time.Second)
	}
}
