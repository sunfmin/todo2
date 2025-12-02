.PHONY: help run test proto clean coverage lint docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Start the application
	go run cmd/api/main.go

test: ## Run all tests
	go test -v ./...

test-race: ## Run tests with race detector
	go test -v -race ./...

coverage: ## Generate test coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

proto: ## Generate protobuf code
	go generate ./...

clean: ## Clean build artifacts and generated files
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf api/gen/v1/*.pb.go
	go clean

build: ## Build the application
	go build -o bin/todo-app cmd/api/main.go

docker-up: ## Start Docker services
	docker-compose up -d

docker-down: ## Stop Docker services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

lint: ## Run linter
	golangci-lint run

tidy: ## Tidy go modules
	go mod tidy

install-tools: ## Install development tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	@echo "Tools installed. Make sure $(go env GOPATH)/bin is in your PATH"