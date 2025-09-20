.PHONY: build test clean client help

# Build the MQTT client CLI
build:
	go build -o bin/mqtt-client ./cmd/client

# Build the client (alias for build)
client: build

# Run tests
test:
	go test ./pkg/mqtt

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run linter (if golangci-lint is installed)
lint:
	golangci-lint run || echo "golangci-lint not installed"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run the client with default settings
run: build
	./bin/mqtt-client

# Help
help:
	@echo "Available targets:"
	@echo "  build    - Build the MQTT client CLI"
	@echo "  client   - Alias for build"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  run      - Build and run the client"
	@echo "  help     - Show this help message"