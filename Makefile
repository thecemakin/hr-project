.PHONY: run build test tidy clean migrate-up migrate-down migrate-force migrate-version swagger setup test-all dc-up dc-down help

# Default app path
APP_ENTRY=./cmd/api/main.go
BINARY_NAME=bin/api

# Database configuration
DB_URL ?= postgres://postgres:postgres@localhost:5432/hr_project?sslmode=disable

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: tidy dc-up migrate-up swagger ## Setup the development environment (deps, db, migrations, swagger)

run: ## Run the application
	go run $(APP_ENTRY)

build: ## Build the application binary
	go build -o $(BINARY_NAME) $(APP_ENTRY)

test: ## Run unit tests
	go test ./...

test-all: ## Run all tests with verbose output
	go test -v ./...

tidy: ## Clean up and verify dependencies
	go mod tidy

dc-up: ## Start docker containers (database)
	docker-compose up -d

dc-down: ## Stop and remove docker containers
	docker-compose down

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	migrate -path ./migrations -database "$(DB_URL)" down

migrate-force: ## Force a specific migration version
	migrate -path ./migrations -database "$(DB_URL)" force $(VERSION)

migrate-version: ## Show current migration version
	migrate -path ./migrations -database "$(DB_URL)" version

clean: ## Clean up built binaries
	rm -rf bin/

swagger: ## Regenerate Swagger/OpenAPI documentation
	@echo "Regenerating Swagger documentation..."
	$(shell go env GOPATH)/bin/swag init -g cmd/api/main.go -d ./ -o docs/openapi --parseDependency --parseInternal
