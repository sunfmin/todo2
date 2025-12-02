# Research: Simple Todo App

**Feature**: 001-todo-app  
**Date**: 2025-12-02  
**Purpose**: Resolve technical unknowns and establish implementation patterns

## Overview

This document resolves the NEEDS CLARIFICATION items from the Technical Context and establishes best practices for implementing a simple todo application following the Go Project Constitution.

## Research Tasks

### 1. Performance Requirements

**Decision**: Target 100 req/s with <200ms p95 latency for MVP

**Rationale**:
- Simple todo app is typically single-user or small team usage
- CRUD operations on small datasets (hundreds of todos per user)
- PostgreSQL can easily handle this load with proper indexing
- Allows room for growth without over-engineering

**Alternatives Considered**:
- 1000 req/s: Over-engineered for MVP, would require caching layer
- 10 req/s: Too conservative, wouldn't validate production readiness

**Implementation Impact**:
- No caching layer needed for MVP
- Standard GORM queries sufficient
- Single database connection pool adequate
- Basic OpenTracing spans sufficient for observability

### 2. Scale Requirements

**Decision**: Support single user with up to 10,000 todos for MVP

**Rationale**:
- Realistic for personal todo app usage
- Tests database performance at reasonable scale
- Pagination becomes important at this scale
- Validates indexing strategy

**Alternatives Considered**:
- Multi-user from start: Adds authentication complexity not in spec
- Unlimited todos: Would require partitioning/archiving strategy
- 100 todos limit: Too restrictive for power users

**Implementation Impact**:
- Add pagination to list endpoint (limit/offset)
- Index on `created_at` for sorting
- No user authentication needed for MVP
- Consider soft deletes for data recovery

### 3. Frontend Technology

**Decision**: Vanilla HTML/CSS/JavaScript (no framework)

**Rationale**:
- Spec says "very simple todo app"
- No build step required
- Easy to understand and modify
- Serves static files from Go server
- Follows KISS principle

**Alternatives Considered**:
- React/Vue: Over-engineered for simple CRUD
- Server-side rendering: Adds complexity, not needed for SPA
- Mobile app: Not specified in requirements

**Implementation Impact**:
- Single `static/index.html` file
- Fetch API for HTTP requests
- No build pipeline needed
- Go server serves static files via `http.FileServer`

### 4. Data Persistence Strategy

**Decision**: PostgreSQL with GORM, no caching layer

**Rationale**:
- Constitution mandates PostgreSQL + GORM
- Performance requirements don't justify caching
- ACID guarantees important for todo operations
- Simplifies architecture

**Alternatives Considered**:
- Redis caching: Premature optimization for MVP
- In-memory only: Violates persistence requirement (FR-005)
- SQLite: Constitution specifies PostgreSQL

**Implementation Impact**:
- Standard GORM CRUD operations
- Database indexes on frequently queried fields
- Transactions for multi-step operations
- testcontainers for integration tests

### 5. API Design Pattern

**Decision**: RESTful JSON API with Protobuf internal types

**Rationale**:
- Constitution requires Protobuf for type safety
- REST is simple and well-understood
- JSON serialization of protobuf structs
- Standard HTTP methods (GET, POST, PUT, DELETE)

**Alternatives Considered**:
- gRPC: Over-engineered for simple web app
- GraphQL: Adds complexity not needed for CRUD
- Plain JSON without protobuf: Violates constitution

**API Endpoints**:
```
POST   /api/v1/todos          # Create todo
GET    /api/v1/todos          # List todos (with pagination)
GET    /api/v1/todos/{id}     # Get single todo
PUT    /api/v1/todos/{id}     # Update todo (toggle complete)
DELETE /api/v1/todos/{id}     # Delete todo
```

**Implementation Impact**:
- Define `todo.proto` with CRUD messages
- Use `r.PathValue()` for path parameters
- Validate with protoc-gen-validate
- Return protobuf structs serialized as JSON

### 6. Error Handling Strategy

**Decision**: Sentinel errors in service layer, HTTP error codes in handler layer

**Rationale**:
- Constitution Principle XIII mandates this pattern
- Clear separation of concerns
- Automatic error mapping via `HandleServiceError()`
- Type-safe with `errors.Is()`

**Service Layer Errors**:
```go
var (
    ErrTodoNotFound     = errors.New("todo not found")
    ErrInvalidInput     = errors.New("invalid input")
    ErrEmptyDescription = errors.New("todo description cannot be empty")
)
```

**HTTP Error Codes**:
```go
var Errors = struct {
    TodoNotFound     ErrorCode  // 404
    InvalidInput     ErrorCode  // 400
    EmptyDescription ErrorCode  // 400
    InternalError    ErrorCode  // 500
}{...}
```

**Implementation Impact**:
- Service methods return sentinel errors
- Handlers use `HandleServiceError()` for automatic mapping
- All errors must have test cases
- Clear error messages for users

### 7. Testing Strategy

**Decision**: Integration tests only, no unit tests or mocks

**Rationale**:
- Constitution Principle I: No mocking
- Real PostgreSQL via testcontainers
- Tests complete HTTP stack via ServeHTTP
- Table-driven design for all test cases

**Test Coverage Requirements**:
- All acceptance scenarios (US1-AS1 through US4-AS3)
- All edge cases from spec
- All error conditions
- >80% code coverage

**Test Structure**:
```go
func TestTodoAPI_Create(t *testing.T) {
    testCases := []struct {
        name     string
        scenario string
        request  *pb.CreateTodoRequest
        wantCode int
        wantErr  string
    }{
        {
            name:     "US1-AS1: Add valid todo",
            scenario: "Given app is open, When user adds 'Buy groceries', Then todo appears",
            request:  &pb.CreateTodoRequest{Description: "Buy groceries"},
            wantCode: http.StatusCreated,
        },
        // ... more test cases
    }
    // ... test execution
}
```

**Implementation Impact**:
- Setup testcontainers in `testutil/db.go`
- Create fixtures in `testutil/fixtures.go`
- Test via root mux ServeHTTP
- Use protocmp for assertions

### 8. Deployment Strategy

**Decision**: Single Docker container with embedded PostgreSQL connection

**Rationale**:
- Simple deployment model
- Environment variables for configuration
- Health check endpoint
- Graceful shutdown

**Configuration**:
```
DATABASE_URL=postgres://user:pass@host:5432/dbname
PORT=8080
LOG_LEVEL=info
TRACING_ENABLED=false  # NoopTracer for MVP
```

**Implementation Impact**:
- Dockerfile with multi-stage build
- Config loaded from environment
- Graceful shutdown on SIGTERM
- Health check at `/health`

## Technology Stack Summary

| Component | Technology | Justification |
|-----------|-----------|---------------|
| Language | Go 1.25+ | Constitution requirement |
| Database | PostgreSQL 15+ | Constitution requirement |
| ORM | GORM | Constitution requirement |
| API Schema | Protocol Buffers | Constitution requirement |
| HTTP Framework | net/http ServeMux | Constitution requirement |
| Tracing | OpenTracing (NoopTracer) | Constitution requirement |
| Testing | testcontainers-go | Constitution requirement |
| Frontend | Vanilla HTML/CSS/JS | Simplicity, no build step |
| Deployment | Docker | Standard containerization |

## Best Practices Applied

### From Constitution

1. **Testing Principles I-IX**: Integration tests, table-driven, edge cases, ServeHTTP, protobuf, fixtures, continuous verification, scenario mapping, coverage
2. **Service Layer (X)**: Public packages, dependency injection, builder pattern, protobuf-only parameters
3. **Tracing (XI)**: OpenTracing spans at HTTP and service level
4. **Context (XII)**: Context-aware operations throughout
5. **Errors (XIII)**: Sentinel errors + HTTP error codes with automatic mapping

### Additional Patterns

1. **Pagination**: Limit/offset for list endpoint
2. **Soft Deletes**: Consider for data recovery (optional)
3. **Timestamps**: `created_at`, `updated_at` for audit trail
4. **Validation**: protoc-gen-validate for input validation
5. **Graceful Shutdown**: Handle SIGTERM properly

## Open Questions

None - All NEEDS CLARIFICATION items resolved.

## Next Steps

Proceed to Phase 1:
1. Generate `data-model.md` with entity definitions
2. Generate `contracts/` with protobuf definitions
3. Generate `quickstart.md` with setup instructions
4. Update agent context with technology choices