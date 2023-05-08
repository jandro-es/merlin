package consumers

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func StartKafkaConsumer() {
	// Set up Kafka consumer configuration
	config := &kafka.ConfigMap{
		// "bootstrap.servers":               "kafka-broker1:9092,kafka-broker2:9092",
		"bootstrap.servers":               "localhost:9092",
		"group.id":                        "my-group",
		"auto.offset.reset":               "earliest",
		"go.application.rebalance.enable": true,
	}

	// Create Kafka consumer client
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Error creating consumer: %v\n", err)
		return
	}
	defer consumer.Close()

	// Subscribe to Kafka topic
	err = consumer.SubscribeTopics([]string{"my-topic"}, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic: %v\n", err)
		return
	}

	// Set up signal handler to handle interrupt signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	// Start consuming Kafka events in a separate goroutine
	go func() {
		for {
			select {
			case sig := <-sigchan:
				// Interrupt signal received, stop consuming events
				fmt.Printf("Caught signal %v: terminating\n", sig)
				return
			default:
				// Poll for Kafka events
				ev := consumer.Poll(100)
				if ev == nil {
					continue
				}
				switch e := ev.(type) {
				case *kafka.Message:
					// Process the Kafka event
					if e.TopicPartition.Error != nil {
						fmt.Printf("Error reading message: %v\n", e.TopicPartition.Error)
					} else if string(e.Key) == "my-event" {
						// My event received, trigger an action
						fmt.Printf("Received my-event: %v\n", string(e.Value))
						triggerAction()
					} else {
						fmt.Printf("Ignoring event: %v\n", string(e.Value))
					}
				case kafka.Error:
					fmt.Printf("Error reading message: %v\n", e)
				default:
					fmt.Printf("Ignored event type: %T\n", e)
				}
			}
		}
	}()

	// Wait for interrupt signal to stop consuming events
	<-sigchan
}

func triggerAction() {
	// Implement your action here
	fmt.Println("My event received, triggering action...")
}
