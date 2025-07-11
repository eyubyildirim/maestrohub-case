# Makefile for the Tower Defence Backend Engine

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=tower-defence-engine
BINARY_UNIX=$(BINARY_NAME)

# --- Build Commands ---

# Build the binary for the current OS
build:
	@echo "Building binary..."
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server

# Build the binary for Linux (for Docker)
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) ./cmd/server

# --- Local Development Commands ---

# Run the backend locally
run:
	@echo "Running backend server..."
	$(GOCMD) run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(BINARY_UNIX)

# --- Docker-related Commands ---

# Build the Docker image for the backend
docker-build: build-linux
	@echo "Building Docker image..."
	docker build -t tower-defence-backend .

# --- Help ---

# Display help information
help:
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  build         Build the application binary"
	@echo "  build-linux   Build for Linux (used for Docker)"
	@echo "  run           Run the application locally"
	@echo "  test          Run unit tests"
	@echo "  clean         Remove build artifacts"
	@echo "  docker-build  Build the Docker image"
	@echo ""

.PHONY: build build-linux run test clean docker-build help
