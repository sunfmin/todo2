# todo2 Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-02

## Active Technologies

### 001-todo-app: Simple Todo App
- **Language**: Go 1.25+
- **Database**: PostgreSQL 15+ with GORM
- **API**: Protocol Buffers (protobuf) with REST/JSON
- **HTTP**: Standard library net/http with ServeMux
- **Tracing**: OpenTracing (NoopTracer for development)
- **Testing**: testcontainers-go with real PostgreSQL
- **Frontend**: Vanilla HTML/CSS/JavaScript (no framework)
- **Validation**: protoc-gen-validate

## Project Structure

```text
api/
├── proto/v1/              # Protobuf definitions
│   └── todo.proto
└── gen/v1/                # Generated protobuf code (PUBLIC)

services/                  # PUBLIC - Business logic
├── todo_service.go
├── errors.go             # Sentinel errors
└── migrations.go         # AutoMigrate() function

handlers/                  # PUBLIC - HTTP handlers
├── todo_handler.go
├── todo_handler_test.go  # Integration tests
├── error_codes.go        # HTTP error codes
└── routes.go             # Routing configuration

internal/                  # INTERNAL - Implementation details
├── models/
│   └── todo.go           # GORM model
├── middleware/
│   ├── logging.go
│   └── tracing.go
└── config/
    └── config.go

cmd/api/                   # Main application
└── main.go

testutil/                  # Test helpers
├── fixtures.go           # Test fixtures
└── db.go                 # Test database setup

static/                    # Frontend files
└── index.html
```

## Commands

```bash
# Development
go run cmd/api/main.go              # Start server
go generate ./...                   # Generate protobuf code

# Testing (requires Docker)
go test -v ./...                    # Run all tests
go test -v -race ./...              # Run with race detector
go test -v -coverprofile=coverage.out ./...  # With coverage

# Database
docker-compose up -d postgres       # Start PostgreSQL
docker exec -it todo-postgres psql -U todouser -d tododb  # Connect to DB

# Protobuf
protoc --go_out=. --go_opt=paths=source_relative \
       --validate_out="lang=go:." --validate_opt=paths=source_relative \
       api/proto/v1/todo.proto
```

## Code Style

- Follow Go Project Constitution (see `.specify/memory/constitution.md`)
- Test-Driven Development (TDD): Write tests before implementation
- Integration tests only (NO mocking), use real PostgreSQL via testcontainers
- Table-driven test design with scenario mapping (US#-AS#)
- Services in PUBLIC packages, models in internal/
- Service methods use protobuf structs only (NO primitives)
- Context-aware operations throughout (ctx as first parameter)
- Sentinel errors in service layer, HTTP error codes in handler layer

## Recent Changes

- 001-todo-app: Added

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
