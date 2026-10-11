.PHONY: build run test clean docker-build docker-run help

# Variables
BINARY_NAME=api-gateway
BUILD_DIR=bin
MAIN_FILE=cmd/server/main.go

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

run: ## Run the application locally
	go run $(MAIN_FILE)

test: ## Run tests
	go test ./... -v -cover

test-race: ## Run tests with race detector
	go test ./... -race -v

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)

docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

docker-compose-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-compose-down: ## Stop all services
	docker-compose down

docker-compose-logs: ## View logs
	docker-compose logs -f api-gateway

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	go fmt ./...

tidy: ## Tidy go modules
	go mod tidy

vet: ## Run go vet
	go vet ./...

check: build test lint fmt vet ## Run all checks

# Development targets
dev-up: docker-compose-up ## Start development environment
dev-down: docker-compose-down ## Stop development environment

# Production targets
prod-build: docker-build ## Build production image
prod-push: ## Push to registry (requires REGISTRY env var)
	docker tag $(BINARY_NAME):latest $(REGISTRY)/$(BINARY_NAME):$(VERSION)
	docker push $(REGISTRY)/$(BINARY_NAME):$(VERSION)
