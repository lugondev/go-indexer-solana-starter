.PHONY: help build run test clean fmt lint docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the indexer binary"
	@echo "  run          - Run the indexer"
	@echo "  test         - Run tests"
	@echo "  test-cover   - Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linters"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"

# Build the binary
build:
	@echo "Building indexer..."
	go build -o bin/indexer cmd/indexer/main.go

# Run the application
run:
	@echo "Running indexer..."
	go run cmd/indexer/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -race ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.txt coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Run linters
lint:
	@echo "Running linters..."
	go vet ./...
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t go-indexer-solana-starter:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --env-file .env -p 8080:8080 go-indexer-solana-starter:latest

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
