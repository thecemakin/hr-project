.PHONY: run build test tidy clean

# Default app path
APP_ENTRY=./cmd/api/main.go
BINARY_NAME=bin/api
run:
	go run ./cmd/api/main.go

build:
	go build -o $(BINARY_NAME) $(APP_ENTRY)

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/
