package main

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lib/pq"
)

func main() {
	// 1. Set up the Kafka Producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "agrinode-kafka:9092"})
	if err != nil {
		log.Fatal("failed to create kafka producer: ", err)
	}
	defer producer.Close()
	topic := "anomaly-events"

	// 2. Set up the PostgreSQL Listener
	connStr := "postgres://agrinode_admin:supersecretpassword@agrinode-postgres:5432/agrinode?sslmode=disable"

	// The listener needs a callback function to report connection drops
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println("listener error: ", err.Error())
		}
	}

	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("anomaly_channel")
	if err != nil {
		log.Fatal("failed to listen to anomaly_channel: ", err)
	}
	defer listener.Close()

	log.Println("event publisher listening to postgres channel 'anomaly_channel'...")

	// 3. The Infinite Loop: Waiting for Broadcasts
	for {
		select {
		case notification := <-listener.Notify:
			if notification == nil {
				continue
			}

			// We received a broadcast! The JSON payload is stored in notification.Extra
			log.Printf("caught anomaly from DB: %s\n", notification.Extra)

			// 4. Forward the JSON payload directly into the Kafka queue
			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(notification.Extra),
			}, nil)

			if err != nil {
				log.Printf("failed to push to kafka: %v\n", err)
			} else {
				log.Println("successfully forwarded anomaly to kafka topic -> anomaly-events")
			}

		case <-time.After(90 * time.Second):
			// Ping the listener every 90 seconds to ensure the connection to Postgres is still alive
			go listener.Ping()
		}
	}
}
