# jusid-charric
Secure messaging backend with MQTT client

A Golang MQTT client for secure messaging applications that can connect to MQTT brokers, subscribe to topics, and publish messages.

## Features

- Connect to MQTT brokers (TCP, SSL/TLS, WebSocket)
- Subscribe to topics with custom message handlers
- Publish messages to topics
- Quality of Service (QoS) support
- Automatic reconnection handling
- Clean and simple API
- Command-line interface for testing

## Installation

```bash
go get github.com/JoshuaDotExe/jusid-charric
```

## Quick Start

### Using the Library

```go
package main

import (
    "log"
    "time"
    
    "github.com/JoshuaDotExe/jusid-charric/pkg/mqtt"
)

func main() {
    // Create configuration
    config := &mqtt.Config{
        BrokerURL:    "tcp://localhost:1883",
        ClientID:     "my-client",
        Username:     "user",
        Password:     "pass",
        CleanSession: true,
        KeepAlive:    30 * time.Second,
    }
    
    // Create client
    client := mqtt.NewClient(config)
    
    // Connect
    if err := client.Connect(); err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect()
    
    // Subscribe
    err := client.Subscribe("my/topic", 1, func(topic string, payload []byte) {
        log.Printf("Received: %s on %s", string(payload), topic)
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Publish
    err = client.PublishString("my/topic", 1, false, "Hello, MQTT!")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Using the CLI

Build and run the CLI client:

```bash
go build -o mqtt-client ./cmd/client
./mqtt-client -broker tcp://localhost:1883 -topic my/messages
```

CLI Commands:
- `pub <topic> <message>` - Publish a message to a topic
- `sub <topic>` - Subscribe to a topic  
- `unsub <topic>` - Unsubscribe from a topic
- `quit` - Exit the application

## Configuration

The MQTT client supports the following configuration options:

- `BrokerURL`: MQTT broker URL (e.g., `tcp://localhost:1883`, `ssl://broker:8883`)
- `ClientID`: Unique client identifier
- `Username`: MQTT username (optional)
- `Password`: MQTT password (optional)
- `CleanSession`: Whether to start a clean session
- `KeepAlive`: Keep-alive interval

## API Reference

### Client Methods

- `NewClient(config *Config) *Client` - Create a new MQTT client
- `Connect() error` - Connect to the MQTT broker
- `Disconnect()` - Disconnect from the broker
- `IsConnected() bool` - Check connection status
- `Subscribe(topic string, qos byte, handler MessageHandler) error` - Subscribe to a topic
- `Unsubscribe(topic string) error` - Unsubscribe from a topic
- `Publish(topic string, qos byte, retained bool, payload []byte) error` - Publish raw bytes
- `PublishString(topic string, qos byte, retained bool, message string) error` - Publish string message

### Message Handler

```go
type MessageHandler func(topic string, payload []byte)
```

## Testing

Run the tests:

```bash
go test ./pkg/mqtt
```

## License

MIT License
