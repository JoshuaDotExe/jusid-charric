package mqtt

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.BrokerURL != "tcp://localhost:1883" {
		t.Errorf("Expected default broker URL to be 'tcp://localhost:1883', got '%s'", config.BrokerURL)
	}

	if config.ClientID != "jusid-charric-client" {
		t.Errorf("Expected default client ID to be 'jusid-charric-client', got '%s'", config.ClientID)
	}

	if !config.CleanSession {
		t.Error("Expected default clean session to be true")
	}

	if config.KeepAlive != 30*time.Second {
		t.Errorf("Expected default keep alive to be 30s, got %v", config.KeepAlive)
	}
}

func TestNewClient(t *testing.T) {
	config := &Config{
		BrokerURL:    "tcp://test:1883",
		ClientID:     "test-client",
		Username:     "testuser",
		Password:     "testpass",
		CleanSession: false,
		KeepAlive:    60 * time.Second,
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.config.BrokerURL != config.BrokerURL {
		t.Errorf("Expected broker URL to be '%s', got '%s'", config.BrokerURL, client.config.BrokerURL)
	}

	if client.config.ClientID != config.ClientID {
		t.Errorf("Expected client ID to be '%s', got '%s'", config.ClientID, client.config.ClientID)
	}
}

func TestNewClientWithNilConfig(t *testing.T) {
	client := NewClient(nil)

	if client == nil {
		t.Fatal("Expected client to be created with default config, got nil")
	}

	// Should use default config
	if client.config.BrokerURL != "tcp://localhost:1883" {
		t.Errorf("Expected default broker URL, got '%s'", client.config.BrokerURL)
	}
}

func TestClientConnectionMethods(t *testing.T) {
	client := NewClient(nil)

	// Test IsConnected when not connected
	if client.IsConnected() {
		t.Error("Expected client to not be connected initially")
	}

	// Note: We don't test actual connection here since it requires a running MQTT broker
	// In a real test environment, you would set up a test broker or use mocks
}
