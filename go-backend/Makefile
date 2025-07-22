# GoNextAuth API Makefile

.PHONY: help dev build run test clean docker-build docker-run swagger deps install-tools

# Default target
help: ## Show this help message
	@echo "GoNextAuth API Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
dev: ## Run the application in development mode with hot reload
	@echo "🚀 Starting development server with Air..."
	air

build: ## Build the application
	@echo "🔨 Building application..."
	go build -o bin/main cmd/server/main.go

run: build ## Build and run the application
	@echo "▶️  Running application..."
	./bin/main

# Dependencies and Tools
deps: ## Download and tidy dependencies
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

install-tools: ## Install development tools (air, swag)
	@echo "🛠️  Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Documentation
swagger: ## Generate Swagger documentation
	@echo "📖 Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o docs
	@echo "✅ Swagger docs generated! Available at: http://localhost:8080/docs/index.html"

# Testing
test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Database
migrate: ## Run database migrations (for development)
	@echo "🗃️  Running database migrations..."
	go run cmd/server/main.go migrate

# Linting and Formatting
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

lint: fmt vet ## Run formatting and vetting

# Cleaning
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/
	rm -rf coverage.out
	rm -rf coverage.html
	rm -rf build-errors.log

# Docker
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t GoNextAuth-api .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env GoNextAuth-api

docker-compose-up: ## Run with Docker Compose
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up --build

docker-compose-down: ## Stop Docker Compose services
	@echo "🐳 Stopping Docker Compose services..."
	docker-compose down

# Production
build-linux: ## Build for Linux (production)
	@echo "🔨 Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux cmd/server/main.go

# Environment setup
setup: install-tools deps swagger ## Complete project setup
	@echo "✅ Project setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Copy .env.example to .env and configure your settings"
	@echo "2. Run 'make dev' to start development server"
	@echo ""

# Health check
health: ## Check if the server is running
	@echo "🏥 Checking server health..."
	@curl -f http://localhost:8080/health || echo "❌ Server is not running"

# Quick development setup
quick-start: setup ## Quick start for new developers
	@echo ""
	@echo "🎉 Welcome to GoNextAuth API!"
	@echo ""
	@echo "📋 Quick Start Checklist:"
	@echo "□ 1. Copy .env.example to .env"
	@echo "□ 2. Configure your database URL in .env"
	@echo "□ 3. Configure your AWS SES credentials in .env"
	@echo "□ 4. Run: make dev"
	@echo ""
	@echo "📖 Documentation will be available at: http://localhost:8080/docs/index.html"
	@echo "🔍 Health check available at: http://localhost:8080/health"
	@echo ""

# Environment validation
check-env: ## Validate environment configuration
	@echo "🔍 Checking environment configuration..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found. Please copy .env.example to .env"; \
		exit 1; \
	fi
	@echo "✅ .env file exists"
	@echo "📋 Environment variables loaded:"
	@echo "  - APP_NAME: $${APP_NAME:-Not set}"
	@echo "  - PORT: $${PORT:-Not set}"
	@echo "  - DATABASE_URL: $${DATABASE_URL:+Set}"
	@echo "  - JWT_SECRET: $${JWT_SECRET:+Set}"
	@echo "  - AWS_SES_FROM_EMAIL: $${AWS_SES_FROM_EMAIL:-Not set}"