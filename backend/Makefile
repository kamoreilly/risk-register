# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@mkdir -p tmp
	@go build -o tmp/main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application (all tests)
test:
	@echo "Testing..."
	@go test ./... -v

# Run only unit tests (skip integration tests)
test-unit:
	@echo "Running unit tests..."
	@go test ./... -v -short

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf tmp

# Helper to check/install AIR
ensure-air:
	@if ! command -v air > /dev/null; then \
		echo "Installing AIR..."; \
		go install github.com/air-verse/air@latest; \
	fi

# Live Reload
watch: ensure-air
	@echo "Watching..."
	@air

.PHONY: all build run test clean watch docker-run docker-down itest ensure-air
