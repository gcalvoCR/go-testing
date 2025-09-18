.PHONY: help build up down restart logs test test-unit test-integration clean run-local deps

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# Docker commands
build: ## Build the Docker image
	docker compose build

up: ## Start the application with PostgreSQL (default)
	docker compose --profile postgres up -d app postgres

up-mongo: ## Start the application with MongoDB
	docker compose --profile mongodb up -d app mongodb

down: ## Stop the application
	docker compose --profile postgres --profile mongodb down

restart: ## Restart the application
	docker compose restart

logs: ## Show application logs
	docker compose logs -f app

logs-mongo: ## Show MongoDB application logs
	docker compose logs -f app

# Testing commands
test: ## Run all tests
	go test ./...

test-unit: ## Run unit tests
	go test ./handlers/...

test-integration: ## Run integration tests (requires server running)
	go test ./tests/...

# Local development
run-local: ## Run the application locally (requires PostgreSQL running)
	go run main.go

run-local-mongo: ## Run the application locally with MongoDB
	MONGODB_URI=mongodb://localhost:27017 DB_TYPE=mongodb go run main.go

deps: ## Download dependencies
	go mod download
	go mod tidy

# Cleanup
clean: ## Clean up Docker containers and volumes
	docker compose --profile postgres --profile mongodb down -v
	docker system prune -f

# Database commands
db-logs: ## Show PostgreSQL logs
	docker compose logs -f postgres

db-shell: ## Connect to PostgreSQL shell
	docker compose exec postgres psql -U postgres -d bankdb

mongo-logs: ## Show MongoDB logs
	docker compose logs -f mongodb

mongo-shell: ## Connect to MongoDB shell
	docker compose exec mongodb mongosh bankdb