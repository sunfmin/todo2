# Simple Todo App

A simple, clean todo application built with Go, PostgreSQL, and vanilla JavaScript.

## Features

- ✅ Add todo items
- ✅ View all todos
- ✅ Mark todos as complete/incomplete
- ✅ Delete todos
- ✅ Persistent storage with PostgreSQL
- ✅ Clean, intuitive interface

## Quick Start

### Prerequisites

- Go 1.25+
- Docker (for PostgreSQL)
- Make (optional but recommended)

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd todo-app
   ```

2. **Start PostgreSQL**
   ```bash
   make docker-up
   # Or: docker-compose up -d
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Generate protobuf code**
   ```bash
   make proto
   # Or: go generate ./...
   ```

5. **Run the application**
   ```bash
   make run
   # Or: go run cmd/api/main.go
   ```

6. **Open your browser**
   ```
   http://localhost:8080
   ```

## Development

### Running Tests

```bash
# Run all tests
make test

# Run with race detector
make test-race

# Generate coverage report
make coverage
```

### Project Structure

```
.
├── api/
│   ├── proto/v1/          # Protobuf definitions
│   └── gen/v1/            # Generated protobuf code
├── services/              # Business logic (public)
├── handlers/              # HTTP handlers (public)
├── internal/
│   ├── models/            # GORM models
│   ├── middleware/        # Middleware
│   └── config/            # Configuration
├── cmd/api/               # Main application
├── testutil/              # Test helpers
├── static/                # Frontend files
└── specs/                 # Documentation
```

### Technology Stack

- **Backend**: Go 1.25+
- **Database**: PostgreSQL 15+
- **ORM**: GORM
- **API**: Protocol Buffers + REST/JSON
- **Testing**: testcontainers-go
- **Frontend**: Vanilla HTML/CSS/JavaScript

### Common Commands

```bash
make help          # Show all available commands
make run           # Start the application
make test          # Run tests
make coverage      # Generate coverage report
make proto         # Generate protobuf code
make clean         # Clean build artifacts
make docker-up     # Start Docker services
make docker-down   # Stop Docker services
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/todos` | Create a new todo |
| GET | `/api/v1/todos` | List all todos (paginated) |
| GET | `/api/v1/todos/{id}` | Get a single todo |
| PUT | `/api/v1/todos/{id}` | Update a todo |
| DELETE | `/api/v1/todos/{id}` | Delete a todo |
| GET | `/health` | Health check |

## Configuration

Set environment variables:

```bash
export DATABASE_URL="postgres://todouser:todopass@localhost:5432/tododb?sslmode=disable"
export PORT=8080
export LOG_LEVEL=info
```

Or create a `.env` file (not tracked in git).

## Testing

The project follows Test-Driven Development (TDD) with integration tests:

- All tests use real PostgreSQL via testcontainers
- Table-driven test design
- >80% code coverage target
- Tests map to acceptance scenarios (US#-AS#)

## Documentation

- [Feature Specification](specs/001-todo-app/spec.md)
- [Implementation Plan](specs/001-todo-app/plan.md)
- [Data Model](specs/001-todo-app/data-model.md)
- [API Contracts](specs/001-todo-app/contracts/)
- [Quickstart Guide](specs/001-todo-app/quickstart.md)
- [Implementation Tasks](specs/001-todo-app/tasks.md)

## Contributing

1. Follow the Go Project Constitution (see `.specify/memory/constitution.md`)
2. Write tests before implementation (TDD)
3. Run tests after every change
4. Ensure >80% code coverage
5. Follow the existing code style

## License

[Your License Here]

## Support

For questions or issues, please refer to the documentation in the `specs/` directory.