# Makefile for go-echo-monolithic

# Variables
PROJECT_NAME := go-echo-monolithic
BUILD_DIR := bin
SERVER_BINARY := $(BUILD_DIR)/server
CLI_BINARY := $(BUILD_DIR)/cli

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# Build flags
LDFLAGS := -w -s

.PHONY: all build server cli test clean deps fmt vet migrate migrate-up migrate-down cleanup run dev help

# Default target
all: build

# Build all binaries
build: server cli

# Build server binary
server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(SERVER_BINARY) ./cmd/server/

# Build CLI binary
cli:
	@echo "Building CLI..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(CLI_BINARY) ./cmd/cli/

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) ./...
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run linting (requires golangci-lint)
lint:
	@echo "Running linting..."
	golangci-lint run

# Database migrations
migrate: migrate-up

migrate-up: cli
	@echo "Running database migrations..."
	$(CLI_BINARY) migrate

migrate-down: cli
	@echo "Rolling back database migrations..."
	$(CLI_BINARY) rollback

# Clean up expired tokens
cleanup: cli
	@echo "Cleaning up expired tokens..."
	$(CLI_BINARY) cleanup

# Run server in development mode
run: server migrate-up
	@echo "Starting server..."
	$(SERVER_BINARY)

# Run server with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	air

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f *.log

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(PROJECT_NAME) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(PROJECT_NAME)

# Show version information
version: cli
	$(CLI_BINARY) version

# Check code quality
check: fmt vet test

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/cosmtrek/air@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build all binaries"
	@echo "  server         - Build server binary"
	@echo "  cli            - Build CLI binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run linting (requires golangci-lint)"
	@echo "  migrate        - Run database migrations"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  cleanup        - Clean up expired tokens"
	@echo "  run            - Build and run server"
	@echo "  dev            - Run with hot reload (requires air)"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  version        - Show version information"
	@echo "  check          - Run format, vet, and test"
	@echo "  install-tools  - Install development tools"
	@echo "  help           - Show this help message"