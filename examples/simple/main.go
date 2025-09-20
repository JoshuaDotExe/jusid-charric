package main

import (
	"log"
	"time"

	"github.com/JoshuaDotExe/jusid-charric/pkg/mqtt"
)

func main() {
	// Create MQTT client with default configuration
	config := mqtt.DefaultConfig()
	config.ClientID = "example-client"
	
	client := mqtt.NewClient(config)

	// Connect to broker (this will fail without a running broker, but demonstrates the API)
	log.Println("Attempting to connect to MQTT broker...")
	if err := client.Connect(); err != nil {
		log.Printf("Failed to connect to MQTT broker: %v", err)
		log.Println("Note: This example requires a running MQTT broker at localhost:1883")
		return
	}
	defer client.Disconnect()

	// Subscribe to a topic
	err := client.Subscribe("example/topic", 1, func(topic string, payload []byte) {
		log.Printf("Received message on %s: %s", topic, string(payload))
	})
	if err != nil {
		log.Printf("Failed to subscribe: %v", err)
		return
	}

	// Publish a message
	message := "Hello from jusid-charric MQTT client!"
	err = client.PublishString("example/topic", 1, false, message)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
		return
	}

	log.Println("Message published successfully!")
	
	// Wait a bit to receive any messages
	time.Sleep(2 * time.Second)
	
	log.Println("Example completed.")
}