package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/JoshuaDotExe/jusid-charric/pkg/mqtt"
)

func main() {
	var (
		brokerURL = flag.String("broker", "tcp://localhost:1883", "MQTT broker URL")
		clientID  = flag.String("client-id", "jusid-charric-client", "MQTT client ID")
		username  = flag.String("username", "", "MQTT username")
		password  = flag.String("password", "", "MQTT password")
		topic     = flag.String("topic", "jusid/messages", "Default topic for messaging")
	)
	flag.Parse()

	// Create MQTT client configuration
	config := &mqtt.Config{
		BrokerURL:    *brokerURL,
		ClientID:     *clientID,
		Username:     *username,
		Password:     *password,
		CleanSession: true,
		KeepAlive:    30 * time.Second,
	}

	// Create and connect MQTT client
	client := mqtt.NewClient(config)
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer client.Disconnect()

	// Subscribe to the default topic
	err := client.Subscribe(*topic, 1, func(topic string, payload []byte) {
		fmt.Printf("\n[RECEIVED] Topic: %s | Message: %s\n> ", topic, string(payload))
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("=== Jusid-Charric MQTT Client ===")
	fmt.Printf("Connected to: %s\n", *brokerURL)
	fmt.Printf("Subscribed to: %s\n", *topic)
	fmt.Println("\nCommands:")
	fmt.Println("  pub <topic> <message>  - Publish a message to a topic")
	fmt.Println("  sub <topic>           - Subscribe to a topic")
	fmt.Println("  unsub <topic>         - Unsubscribe from a topic")
	fmt.Println("  quit                  - Exit the application")
	fmt.Println("\nEnter commands or type messages to publish to the default topic:")

	// Start command input goroutine
	inputChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if scanner.Scan() {
				inputChan <- scanner.Text()
			}
		}
	}()

	// Main event loop
	for {
		select {
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return

		case input := <-inputChan:
			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}

			parts := strings.SplitN(input, " ", 3)
			command := strings.ToLower(parts[0])

			switch command {
			case "quit", "exit", "q":
				fmt.Println("Goodbye!")
				return

			case "pub", "publish":
				if len(parts) < 3 {
					fmt.Println("Usage: pub <topic> <message>")
					continue
				}
				pubTopic := parts[1]
				message := parts[2]

				if err := client.PublishString(pubTopic, 1, false, message); err != nil {
					fmt.Printf("Failed to publish message: %v\n", err)
				} else {
					fmt.Printf("Published to %s: %s\n", pubTopic, message)
				}

			case "sub", "subscribe":
				if len(parts) < 2 {
					fmt.Println("Usage: sub <topic>")
					continue
				}
				subTopic := parts[1]

				err := client.Subscribe(subTopic, 1, func(topic string, payload []byte) {
					fmt.Printf("\n[RECEIVED] Topic: %s | Message: %s\n> ", topic, string(payload))
				})
				if err != nil {
					fmt.Printf("Failed to subscribe to %s: %v\n", subTopic, err)
				}

			case "unsub", "unsubscribe":
				if len(parts) < 2 {
					fmt.Println("Usage: unsub <topic>")
					continue
				}
				unsubTopic := parts[1]

				if err := client.Unsubscribe(unsubTopic); err != nil {
					fmt.Printf("Failed to unsubscribe from %s: %v\n", unsubTopic, err)
				}

			case "help", "h":
				fmt.Println("\nCommands:")
				fmt.Println("  pub <topic> <message>  - Publish a message to a topic")
				fmt.Println("  sub <topic>           - Subscribe to a topic")
				fmt.Println("  unsub <topic>         - Unsubscribe from a topic")
				fmt.Println("  quit                  - Exit the application")

			default:
				// Treat any other input as a message to publish to the default topic
				if err := client.PublishString(*topic, 1, false, input); err != nil {
					fmt.Printf("Failed to publish message: %v\n", err)
				} else {
					fmt.Printf("Published to %s: %s\n", *topic, input)
				}
			}
		}
	}
}
