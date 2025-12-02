# Implementation Plan: Simple Todo App

**Branch**: `001-todo-app` | **Date**: 2025-12-02 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-todo-app/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

A simple todo application that allows users to add, view, mark as complete/incomplete, and delete todo items. The app must persist data across sessions and provide a clean, intuitive interface for task management. This is a web-based application using Go backend with PostgreSQL for data persistence.

## Technical Context

**Stack**: Go 1.25+, PostgreSQL 15+, GORM, Protobuf, OpenTracing, testcontainers-go
**Project Type**: Single web application (API + static frontend)
**Target**: Linux server (containerized)
**Performance**: 100 req/s with <200ms p95 latency (resolved in research.md)
**Scale**: Single user with up to 10,000 todos (resolved in research.md)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Testing Requirements
- ✅ Integration tests with real PostgreSQL (Principle I)
- ✅ Table-driven test design (Principle II)
- ✅ Comprehensive edge case coverage (Principle III)
- ✅ ServeHTTP endpoint testing via root mux (Principle IV)
- ✅ Protobuf data structures for API (Principle V)
- ✅ Continuous test verification after every change (Principle VI)
- ✅ Root cause tracing for debugging (Principle VII)
- ✅ Acceptance scenario coverage - All US#-AS# scenarios mapped to tests (Principle VIII)
- ✅ Test coverage >80% for business logic (Principle IX)

### Architecture Requirements
- ✅ Service layer with dependency injection (Principle X)
- ✅ Services in public packages, models in internal/ (Principle X)
- ✅ Service methods use protobuf structs only (Principle X)
- ✅ Distributed tracing with OpenTracing (Principle XI)
- ✅ Context-aware operations throughout (Principle XII)
- ✅ Comprehensive error handling with sentinel errors (Principle XIII)

### Violations/Justifications
None - This is a standard CRUD application that fits the constitution perfectly.

### Post-Design Re-evaluation (Phase 1 Complete)

**Date**: 2025-12-02

All constitution requirements remain satisfied after design phase:

✅ **Testing**: All 13 acceptance scenarios (US1-AS1 through US4-AS3) will be mapped to integration tests using table-driven design with real PostgreSQL via testcontainers.

✅ **Architecture**:
- Services in public `services/` package returning protobuf types
- Models in `internal/models/` (GORM only)
- Handlers in public `handlers/` package with integration tests
- Service methods use protobuf structs only (no primitives)

✅ **Error Handling**:
- Service layer: Sentinel errors (ErrTodoNotFound, ErrInvalidInput, ErrEmptyDescription)
- HTTP layer: Error code singleton with automatic ServiceErr mapping
- All errors will have test cases

✅ **Tracing**: OpenTracing spans at HTTP and service level (NoopTracer for development)

✅ **Context**: All operations context-aware (ctx as first parameter)

✅ **Data Model**: Single entity (Todo) with proper validation, indexing, and protobuf contracts defined

**Conclusion**: Design is constitution-compliant. Ready for Phase 2 (task breakdown).

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., cmd/api, cmd/worker). The delivered plan must not include
  Option labels.
-->

```text
# [REMOVE IF UNUSED] Option 1: Single Go API (DEFAULT)
# Services MUST be in public packages (not internal/) for cross-app reusability
api/
├── proto/                  # Protobuf definitions (.proto files)
│   └── v1/
│       └── *.proto
└── gen/                    # Protobuf generated code (PUBLIC - must be importable)
    └── v1/
        ├── *.pb.go
        └── *.pb.validate.go

services/                   # PUBLIC package - business logic (reusable by external apps)
├── product_service.go
├── errors.go              # Sentinel errors (ErrNotFound, ErrDuplicateSKU, etc.)
└── migrations.go          # AutoMigrate() function for external apps

handlers/                   # PUBLIC package - HTTP handlers (reusable)
├── product_handler.go
├── product_handler_test.go
├── error_codes.go         # HTTP error code singleton with ServiceErr mapping
└── routes.go              # Shared routing configuration

internal/                   # INTERNAL - implementation details only
├── models/                # GORM models (internal - services return protobuf)
│   └── product.go
├── middleware/            # App-specific middleware (logging, CORS)
│   ├── logging.go
│   └── tracing.go
└── config/                # Configuration loading

cmd/
└── api/                   # Main application entry point
    └── main.go

testutil/                   # Test helpers and fixtures
├── fixtures.go            # CreateTestXxx() functions with default values
└── db.go                  # setupTestDB() with testcontainers

# [REMOVE IF UNUSED] Option 2: Multiple Go services (microservices)
# Each service follows Option 1 structure
service-a/
├── api/
│   ├── proto/            # Protobuf definitions
│   └── gen/              # Generated code
├── services/              # PUBLIC - reusable
├── handlers/              # HTTP handlers (with *_test.go integration tests)
├── internal/              # Implementation details
├── testutil/              # Test helpers and fixtures
└── cmd/

service-b/
├── api/
│   ├── proto/            # Protobuf definitions
│   └── gen/              # Generated code
├── services/              # PUBLIC - reusable
├── handlers/              # HTTP handlers (with *_test.go integration tests)
├── internal/              # Implementation details
├── testutil/              # Test helpers and fixtures
└── cmd/

shared/                    # Shared libraries across services
├── tracing/
├── logging/
└── database/
```

**Structure Decision**: Using Option 1 (Single Go API) as this is a simple todo application that doesn't require microservices architecture. The structure will be:

```text
api/
├── proto/v1/
│   └── todo.proto          # Todo API definitions
└── gen/v1/
    ├── todo.pb.go
    └── todo.pb.validate.go

services/
├── todo_service.go         # Business logic
├── errors.go              # Sentinel errors
└── migrations.go          # AutoMigrate() function

handlers/
├── todo_handler.go        # HTTP handlers
├── todo_handler_test.go   # Integration tests
├── error_codes.go         # HTTP error codes
└── routes.go              # Routing configuration

internal/
├── models/
│   └── todo.go            # GORM model
├── middleware/
│   ├── logging.go
│   └── tracing.go
└── config/
    └── config.go

cmd/api/
└── main.go                # Application entry point

testutil/
├── fixtures.go            # Test fixtures
└── db.go                  # Test database setup

static/                    # Frontend files (HTML/CSS/JS)
└── index.html
```

**Architecture** (Constitution Principle X):
- Services/handlers: PUBLIC packages (return protobuf, reusable)
- Models: `internal/models/` (GORM only, never exposed)
- Protobuf: PUBLIC `api/gen/` (external apps need these)
- `AutoMigrate()`: Exported in `services/migrations.go`

## Testing Strategy

### Test-First Development (TDD)

TDD workflow (Constitution Development Workflow):

1. **Design**: Define API in `.proto` files → generate code
2. **Red**: Write integration tests → verify FAIL
3. **Green**: Implement → run tests → verify PASS
4. **Refactor**: Improve code → run tests after each change
5. **Complete**: Done only when ALL tests pass

### Integration Testing Requirements

Constitution Testing Principles I-IX:

- **Integration tests ONLY** (NO mocking), real PostgreSQL via testcontainers
- **Table-driven** with `name` fields
- **Edge cases MANDATORY**: Input validation, boundaries, auth, data state, database, HTTP
- **ServeHTTP testing** via root mux (NOT individual handlers)
- **Protobuf** structs with `protocmp` assertions
- **Derive from fixtures** (NOT response, except UUIDs/timestamps)
- **Run tests** after EVERY change (Principle VI)
- **Map scenarios** to tests (US#-AS#, Principle VIII)
- **Coverage >80%** (Principle IX)

### Error Handling Strategy

Constitution Principle XI:

- **Service**: Sentinel errors (`var ErrXxx`), wrap with `fmt.Errorf("%w")`
- **HTTP**: Error singleton with `ServiceErr` mapping, automatic via `HandleServiceError()`
- **Testing**: ALL errors must have test cases

### Test Database Isolation

- **testcontainers-go** with PostgreSQL (Docker required)
- **Truncation**: `defer truncateTables(db, "tables...")` with CASCADE
- **Parallel**: `t.Parallel()` safe

### Context-Aware Operations

Constitution Principle XII: Services accept `context.Context`, handlers use `r.Context()`, database uses `db.WithContext(ctx)`, tests verify cancellation.

### Distributed Tracing

Constitution Principle XI: HTTP endpoints create OpenTracing spans, services create child spans, database as ONE span per transaction (NOT per query). Dev uses `NoopTracer{}`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
