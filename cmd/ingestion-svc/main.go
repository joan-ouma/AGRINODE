package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("  Agri-Node Ingestion Service Starting  ")
	fmt.Println("========================================")

	// Note for future use:
	// This is where we will eventually initialize the Kafka adapter
	// and inject it into our IngestionUseCase, like this:
	//
	// kafkaPublisher := adapters.NewKafkaPublisher(...)
	// ingestionUseCase := usecases.NewIngestionUseCase(kafkaPublisher)
	//
	// Then we will connect to the MQTT broker and start passing data to the use case.

	log.Println("Service boilerplate initialized successfully. Ready for Week 2!")
}
