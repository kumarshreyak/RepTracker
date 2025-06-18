# GymLog Backend Makefile

.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-down deploy-local deploy-prod

# Default target
help:
	@echo "GymLog Backend - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  setup          - Set up development environment"
	@echo "  run            - Run the server locally"
	@echo "  test           - Run tests"
	@echo "  build          - Build the binary"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container locally"
	@echo "  docker-compose-up   - Start with Docker Compose"
	@echo "  docker-compose-down - Stop Docker Compose"
	@echo ""
	@echo "Deployment:"
	@echo "  deploy-local   - Deploy to local Docker"
	@echo "  deploy-prod    - Deploy to Google Cloud Run"
	@echo ""

# Variables
BINARY_NAME=server
DOCKER_IMAGE=gymlog-backend
GCP_PROJECT_ID?=your-gcp-project-id

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		cp env.template .env; \
		echo "Created .env file from template. Please update with your values."; \
	else \
		echo ".env file already exists."; \
	fi
	@go mod tidy
	@echo "Setup complete!"

# Build the binary
build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) ./cmd/server
	@echo "Binary built: $(BINARY_NAME)"

# Run the server locally
run:
	@echo "Starting server..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f main
	@go clean

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Docker Compose commands
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-compose-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

docker-compose-logs:
	@docker-compose logs -f

# Deployment commands
deploy-local: docker-build
	@echo "Deploying locally with Docker..."
	@docker run -d -p 8080:8080 --name gymlog-backend --env-file .env $(DOCKER_IMAGE)
	@echo "Local deployment complete! API available at http://localhost:8080"

deploy-prod:
	@echo "Deploying to Google Cloud Run..."
	@if [ "$(GCP_PROJECT_ID)" = "your-gcp-project-id" ]; then \
		echo "Error: Please set GCP_PROJECT_ID to your actual project ID"; \
		echo "Usage: make deploy-prod GCP_PROJECT_ID=your-actual-project-id"; \
		exit 1; \
	fi
	@./deploy.sh production $(GCP_PROJECT_ID)

# Development helpers
dev-setup: setup
	@echo "Setting up for development..."
	@docker-compose up -d mongodb
	@echo "MongoDB started. You can now run 'make run' to start the server."

dev-logs:
	@docker-compose logs -f mongodb

# Health check
health:
	@echo "Checking service health..."
	@curl -s http://localhost:8080/health | jq . || echo "Service not available or jq not installed"

# Database operations
populate-exercises:
	@echo "Populating exercises database..."
	@cd scripts && go run populate_exercises.go

# Git hooks
install-hooks:
	@echo "Installing git hooks..."
	@echo '#!/bin/sh\nmake test' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hook installed" 