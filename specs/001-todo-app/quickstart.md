# Quickstart Guide: Simple Todo App

**Feature**: 001-todo-app  
**Date**: 2025-12-02  
**Purpose**: Get the todo app running locally in under 10 minutes

## Prerequisites

Before you begin, ensure you have:

- **Go 1.25+**: [Download](https://go.dev/dl/)
- **Docker**: [Download](https://www.docker.com/products/docker-desktop/) (for PostgreSQL and tests)
- **Git**: For version control
- **Make**: Optional but recommended

Verify installations:
```bash
go version        # Should show 1.25 or higher
docker --version  # Should show Docker version
git --version     # Should show Git version
```

## Quick Start (5 minutes)

### 1. Clone and Setup

```bash
# Clone repository
git clone <repository-url>
cd todo-app

# Checkout feature branch
git checkout 001-todo-app

# Install dependencies
go mod download
```

### 2. Start PostgreSQL

```bash
# Using Docker Compose (recommended)
docker-compose up -d postgres

# Or using Docker directly
docker run -d \
  --name todo-postgres \
  -e POSTGRES_USER=todouser \
  -e POSTGRES_PASSWORD=todopass \
  -e POSTGRES_DB=tododb \
  -p 5432:5432 \
  postgres:15-alpine
```

### 3. Generate Protobuf Code

```bash
# Install protoc plugins (first time only)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/envoyproxy/protoc-gen-validate@latest

# Generate code
go generate ./...
```

### 4. Run Database Migrations

```bash
# Migrations run automatically on startup
# Or run manually:
go run cmd/migrate/main.go
```

### 5. Start the Server

```bash
# Development mode
go run cmd/api/main.go

# Or using Make
make run
```

Server starts at: http://localhost:8080

### 6. Open the App

Open your browser to: http://localhost:8080

You should see the todo app interface!

## Development Workflow

### Project Structure

```
todo-app/
├── api/
│   ├── proto/v1/          # Protobuf definitions
│   │   └── todo.proto
│   └── gen/v1/            # Generated protobuf code
├── services/              # Business logic (PUBLIC)
│   ├── todo_service.go
│   ├── errors.go
│   └── migrations.go
├── handlers/              # HTTP handlers (PUBLIC)
│   ├── todo_handler.go
│   ├── todo_handler_test.go
│   ├── error_codes.go
│   └── routes.go
├── internal/              # Internal implementation
│   ├── models/
│   │   └── todo.go
│   ├── middleware/
│   └── config/
├── cmd/api/               # Main application
│   └── main.go
├── testutil/              # Test helpers
│   ├── fixtures.go
│   └── db.go
├── static/                # Frontend files
│   └── index.html
└── specs/001-todo-app/    # Documentation
    ├── spec.md
    ├── plan.md
    ├── research.md
    ├── data-model.md
    ├── quickstart.md
    └── contracts/
```

### Test-Driven Development (TDD)

**CRITICAL**: Follow TDD workflow from constitution:

1. **Design**: Define API in `.proto` files
2. **Red**: Write failing integration test
3. **Green**: Implement to make test pass
4. **Refactor**: Improve code, keep tests green
5. **Complete**: Done when ALL tests pass

### Running Tests

```bash
# Run all tests (requires Docker)
go test -v ./...

# Run with race detector
go test -v -race ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestTodoAPI_Create ./handlers

# Run tests in watch mode (using entr)
find . -name "*.go" | entr -c go test -v ./...
```

**Important**: Tests use testcontainers and require Docker to be running.

### Making Changes

#### 1. Update Protobuf Schema

```bash
# Edit api/proto/v1/todo.proto
vim api/proto/v1/todo.proto

# Regenerate code
go generate ./...

# Verify compilation
go build ./...
```

#### 2. Write Integration Test

```bash
# Edit handlers/todo_handler_test.go
vim handlers/todo_handler_test.go

# Run test (should FAIL)
go test -v ./handlers -run TestTodoAPI_YourNewTest
```

#### 3. Implement Feature

```bash
# Update service layer
vim services/todo_service.go

# Update handler layer
vim handlers/todo_handler.go

# Run test (should PASS)
go test -v ./handlers -run TestTodoAPI_YourNewTest
```

#### 4. Run Full Test Suite

```bash
# All tests must pass before commit
go test -v ./...

# Check coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

### Code Generation

```bash
# Generate protobuf code
go generate ./...

# Generate mocks (if needed)
mockgen -source=services/todo_service.go -destination=mocks/todo_service_mock.go

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

### Database Operations

```bash
# Connect to database
docker exec -it todo-postgres psql -U todouser -d tododb

# View todos
SELECT * FROM todos ORDER BY created_at DESC;

# Reset database
docker-compose down -v
docker-compose up -d postgres
go run cmd/migrate/main.go
```

## API Testing

### Using curl

```bash
# Create todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"description": "Buy groceries"}'

# List todos
curl http://localhost:8080/api/v1/todos

# Get single todo
curl http://localhost:8080/api/v1/todos/{id}

# Update todo
curl -X PUT http://localhost:8080/api/v1/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'

# Delete todo
curl -X DELETE http://localhost:8080/api/v1/todos/{id}

# Health check
curl http://localhost:8080/health
```

### Using HTTPie (prettier output)

```bash
# Install HTTPie
brew install httpie  # macOS
apt install httpie   # Ubuntu

# Create todo
http POST localhost:8080/api/v1/todos description="Buy groceries"

# List todos
http GET localhost:8080/api/v1/todos

# Update todo
http PUT localhost:8080/api/v1/todos/{id} completed:=true
```

### Using Postman

1. Import `specs/001-todo-app/contracts/openapi.yaml`
2. Set base URL to `http://localhost:8080/api/v1`
3. Run requests from collection

## Configuration

### Environment Variables

```bash
# Database
export DATABASE_URL="postgres://todouser:todopass@localhost:5432/tododb?sslmode=disable"

# Server
export PORT=8080
export LOG_LEVEL=info

# Tracing (optional)
export TRACING_ENABLED=false
export JAEGER_ENDPOINT=http://localhost:14268/api/traces
```

### Configuration File

Create `.env` file:
```env
DATABASE_URL=postgres://todouser:todopass@localhost:5432/tododb?sslmode=disable
PORT=8080
LOG_LEVEL=debug
TRACING_ENABLED=false
```

Load with:
```bash
# Using godotenv
go get github.com/joho/godotenv
```

## Troubleshooting

### Tests Failing

**Problem**: Tests fail with "cannot connect to Docker daemon"

**Solution**:
```bash
# Start Docker Desktop
# Or start Docker daemon
sudo systemctl start docker
```

**Problem**: Tests fail with "port already in use"

**Solution**:
```bash
# Kill process using port
lsof -ti:5432 | xargs kill -9

# Or use different port
docker run -p 5433:5432 postgres:15-alpine
```

### Database Issues

**Problem**: "relation does not exist"

**Solution**:
```bash
# Run migrations
go run cmd/migrate/main.go

# Or restart with fresh database
docker-compose down -v
docker-compose up -d
```

**Problem**: "connection refused"

**Solution**:
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection string
echo $DATABASE_URL
```

### Protobuf Issues

**Problem**: "protoc-gen-go: program not found"

**Solution**:
```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/envoyproxy/protoc-gen-validate@latest

# Add to PATH
export PATH="$PATH:$(go env GOPATH)/bin"
```

**Problem**: "cannot find package"

**Solution**:
```bash
# Regenerate protobuf code
go generate ./...

# Download dependencies
go mod download
```

## Next Steps

### Development Tasks

1. **Review Spec**: Read `specs/001-todo-app/spec.md` for requirements
2. **Review Plan**: Read `specs/001-todo-app/plan.md` for architecture
3. **Review Contracts**: Read `specs/001-todo-app/contracts/` for API details
4. **Write Tests**: Start with acceptance scenarios (US#-AS#)
5. **Implement Features**: Follow TDD workflow

### Production Deployment

See deployment guide (coming soon) for:
- Docker image building
- Kubernetes deployment
- Environment configuration
- Monitoring setup
- Backup strategy

## Resources

### Documentation

- [Feature Spec](spec.md) - Requirements and acceptance criteria
- [Implementation Plan](plan.md) - Architecture and design decisions
- [Research](research.md) - Technical decisions and rationale
- [Data Model](data-model.md) - Entity definitions and schema
- [API Contracts](contracts/) - Protobuf and OpenAPI specs

### External Resources

- [Go Project Constitution](../../.specify/memory/constitution.md) - Development principles
- [Protocol Buffers Guide](https://protobuf.dev/getting-started/gotutorial/)
- [GORM Documentation](https://gorm.io/docs/)
- [testcontainers-go](https://golang.testcontainers.org/)

## Support

For questions or issues:
1. Check troubleshooting section above
2. Review constitution for development principles
3. Check existing tests for examples
4. Ask team for help

## Quick Reference

```bash
# Common commands
make run          # Start server
make test         # Run tests
make coverage     # Generate coverage report
make proto        # Generate protobuf code
make lint         # Run linter
make clean        # Clean build artifacts

# Docker commands
docker-compose up -d      # Start services
docker-compose down       # Stop services
docker-compose logs -f    # View logs

# Database commands
docker exec -it todo-postgres psql -U todouser -d tododb
\dt                       # List tables
\d todos                  # Describe todos table
SELECT * FROM todos;      # Query todos
```

Happy coding! 🚀