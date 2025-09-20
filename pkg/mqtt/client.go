package mqtt

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client represents an MQTT client for secure messaging
type Client struct {
	client mqtt.Client
	config *Config
}

// Config holds the configuration for the MQTT client
type Config struct {
	BrokerURL    string
	ClientID     string
	Username     string
	Password     string
	CleanSession bool
	KeepAlive    time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		BrokerURL:    "tcp://localhost:1883",
		ClientID:     "jusid-charric-client",
		CleanSession: true,
		KeepAlive:    30 * time.Second,
	}
}

// MessageHandler defines the function signature for message handlers
type MessageHandler func(topic string, payload []byte)

// NewClient creates a new MQTT client with the given configuration
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.BrokerURL)
	opts.SetClientID(config.ClientID)
	opts.SetCleanSession(config.CleanSession)
	opts.SetKeepAlive(config.KeepAlive)

	if config.Username != "" {
		opts.SetUsername(config.Username)
	}
	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	// Set default connection lost handler
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	})

	// Set default reconnecting handler
	opts.SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
		log.Println("Attempting to reconnect...")
	})

	client := mqtt.NewClient(opts)

	return &Client{
		client: client,
		config: config,
	}
}

// Connect establishes a connection to the MQTT broker
func (c *Client) Connect() error {
	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	log.Printf("Connected to MQTT broker at %s", c.config.BrokerURL)
	return nil
}

// Disconnect closes the connection to the MQTT broker
func (c *Client) Disconnect() {
	c.client.Disconnect(250)
	log.Println("Disconnected from MQTT broker")
}

// IsConnected returns true if the client is connected to the broker
func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}

// Subscribe subscribes to a topic with the given Quality of Service (QoS) level
func (c *Client) Subscribe(topic string, qos byte, handler MessageHandler) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		if handler != nil {
			handler(msg.Topic(), msg.Payload())
		} else {
			log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
		}
	}

	token := c.client.Subscribe(topic, qos, messageHandler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	log.Printf("Subscribed to topic: %s (QoS: %d)", topic, qos)
	return nil
}

// Unsubscribe unsubscribes from a topic
func (c *Client) Unsubscribe(topic string) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	token := c.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topic %s: %w", topic, token.Error())
	}

	log.Printf("Unsubscribed from topic: %s", topic)
	return nil
}

// Publish publishes a message to a topic with the given QoS level
func (c *Client) Publish(topic string, qos byte, retained bool, payload []byte) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	token := c.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, token.Error())
	}

	log.Printf("Published message to topic: %s", topic)
	return nil
}

// PublishString publishes a string message to a topic
func (c *Client) PublishString(topic string, qos byte, retained bool, message string) error {
	return c.Publish(topic, qos, retained, []byte(message))
}
