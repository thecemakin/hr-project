.PHONY: run build test tidy clean migrate-up migrate-down

# Default app path
APP_ENTRY=./cmd/api/main.go
BINARY_NAME=bin/api

# Database configuration
DB_URL ?= postgres://postgres:postgres@database.hr-project.orb.local:5432/hr_project?sslmode=disable

run:
	go run ./cmd/api/main.go

build:
	go build -o $(BINARY_NAME) $(APP_ENTRY)

test:
	go test ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path ./migrations -database $(DB_URL) up

migrate-down:
	migrate -path ./migrations -database $(DB_URL) down

migrate-force:
	migrate -path ./migrations -database $(DB_URL) force $(VERSION)

migrate-version:
	migrate -path ./migrations -database $(DB_URL) version

clean:
	rm -rf bin/
