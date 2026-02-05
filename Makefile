.PHONY: help build run test docker-build docker-up docker-down docker-logs migrate-up migrate-down lint fmt deps clean

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker Compose services"
	@echo "  make docker-down   - Stop Docker Compose services"
	@echo "  make docker-logs   - View Docker Compose logs"
	@echo "  make migrate-up    - Apply database migrations"
	@echo "  make migrate-down  - Rollback database migrations"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"
	@echo "  make deps          - Download dependencies"
	@echo "  make clean         - Clean build artifacts"

build:
	go build -o bin/taskmanager cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

docker-build:
	docker build -t taskmanager:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

migrate-up:
	docker exec -i taskmanager_postgres psql -U taskmanager -d taskmanager < internal/database/migrations/001_init_schema.up.sql

migrate-down:
	docker exec -i taskmanager_postgres psql -U taskmanager -d taskmanager < internal/database/migrations/001_init_schema.down.sql

lint:
	golangci-lint run

fmt:
	go fmt ./...

deps:
	go mod download
	go mod tidy

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
