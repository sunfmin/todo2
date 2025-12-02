# Implementation Tasks: Simple Todo App

**Feature**: 001-todo-app  
**Branch**: `001-todo-app`  
**Date**: 2025-12-02  
**Status**: Ready for Implementation

## Overview

This document breaks down the implementation into actionable tasks organized by user story priority. Each phase represents a complete, independently testable increment following Test-Driven Development (TDD).

**Total Tasks**: 45  
**Estimated Duration**: 3-5 days  
**MVP Scope**: Phase 3 (User Story 1 + 4) - Core todo functionality

## Task Format

```
- [ ] [TaskID] [P] [Story] Description with file path
```

- **TaskID**: Sequential number (T001, T002...)
- **[P]**: Parallelizable (can be done simultaneously with other [P] tasks)
- **[Story]**: User story label (US1, US2, US3, US4)
- **Description**: Clear action with exact file path

## Implementation Strategy

### MVP-First Approach
1. **Phase 1-2**: Setup infrastructure (blocking prerequisites)
2. **Phase 3**: Implement US1 + US4 (P1 stories) = **MVP Release**
3. **Phase 4**: Add US2 (P2) = Enhanced functionality
4. **Phase 5**: Add US3 (P3) = Complete feature set
5. **Phase 6**: Polish and optimization

### TDD Workflow (Per Task)
1. Write failing integration test
2. Implement minimal code to pass
3. Run full test suite
4. Refactor if needed
5. Verify all tests still pass

---

## Phase 1: Project Setup

**Goal**: Initialize project structure and dependencies

**Duration**: 30 minutes

### Tasks

- [ ] T001 Initialize Go module with `go mod init github.com/yourorg/todo-app`
- [ ] T002 Create project directory structure per plan.md
- [ ] T003 [P] Create go.mod with dependencies (gorm.io/gorm, gorm.io/driver/postgres, google.golang.org/protobuf, github.com/opentracing/opentracing-go, github.com/testcontainers/testcontainers-go)
- [ ] T004 [P] Create .gitignore file (api/gen/, *.pb.go, coverage.out, bin/)
- [ ] T005 [P] Create docker-compose.yml with PostgreSQL service
- [ ] T006 [P] Create Makefile with common commands (run, test, proto, clean)
- [ ] T007 [P] Create README.md with quickstart instructions

---

## Phase 2: Foundational Layer

**Goal**: Set up core infrastructure needed by all user stories

**Duration**: 2-3 hours

**Blocking**: Must complete before any user story implementation

### Protobuf & Code Generation

- [ ] T008 Copy api/proto/v1/todo.proto from specs/001-todo-app/contracts/todo.proto
- [ ] T009 Add //go:generate directive to api/proto/v1/todo.proto
- [ ] T010 Run `go generate ./...` to generate protobuf code in api/gen/v1/
- [ ] T011 Verify generated files compile with `go build ./...`

### Database Layer

- [ ] T012 [P] Create internal/models/todo.go with GORM Todo model
- [ ] T013 [P] Create services/migrations.go with AutoMigrate() function
- [ ] T014 [P] Create internal/config/config.go for database configuration

### Service Layer Foundation

- [ ] T015 Create services/errors.go with sentinel errors (ErrTodoNotFound, ErrInvalidInput, ErrEmptyDescription)
- [ ] T016 Create services/todo_service.go with TodoService interface and builder pattern
- [ ] T017 Add helper functions toProto(), fromCreateRequest(), applyUpdate() in services/todo_service.go

### HTTP Layer Foundation

- [ ] T018 Create handlers/error_codes.go with error code singleton and ServiceErr mapping
- [ ] T019 Create handlers/todo_handler.go with TodoHandler struct
- [ ] T020 Create handlers/routes.go with SetupRoutes() function
- [ ] T021 [P] Create internal/middleware/logging.go for request logging
- [ ] T022 [P] Create internal/middleware/tracing.go with OpenTracing NoopTracer

### Test Infrastructure

- [ ] T023 Create testutil/db.go with setupTestDB() using testcontainers
- [ ] T024 Create testutil/fixtures.go with CreateTestTodo() helper
- [ ] T025 Create handlers/todo_handler_test.go with test setup boilerplate

### Main Application

- [ ] T026 Create cmd/api/main.go with server initialization, database connection, graceful shutdown
- [ ] T027 Add /health endpoint in handlers/routes.go
- [ ] T028 Test server starts successfully with `go run cmd/api/main.go`

---

## Phase 3: User Story 1 + 4 (P1) - MVP Core

**Goal**: Users can add and view todos (core functionality)

**Duration**: 4-6 hours

**Independent Test**: User can add "Buy groceries", see it in the list, refresh page, and todo persists

**Value Delivered**: Functional todo app with persistence

### User Story 1: Add Todo Items

#### Tests First (TDD Red Phase)

- [ ] T029 [US1] Write test US1-AS1 in handlers/todo_handler_test.go: Add "Buy groceries" → appears in list
- [ ] T030 [US1] Write test US1-AS2 in handlers/todo_handler_test.go: Add todo to existing list → both visible
- [ ] T031 [US1] Write test US1-AS3 in handlers/todo_handler_test.go: Empty todo → validation error
- [ ] T032 [US1] Write edge case tests: whitespace-only, 500+ chars, special characters, emojis
- [ ] T033 [US1] Run tests with `go test -v ./handlers` → verify all FAIL

#### Implementation (TDD Green Phase)

- [ ] T034 [US1] Implement TodoService.Create() in services/todo_service.go with validation
- [ ] T035 [US1] Implement TodoHandler.Create() in handlers/todo_handler.go (POST /api/v1/todos)
- [ ] T036 [US1] Register Create route in handlers/routes.go
- [ ] T037 [US1] Run tests with `go test -v ./handlers` → verify all PASS
- [ ] T038 [US1] Run full test suite with `go test -v ./...` → verify no regressions

### User Story 4: View All Todos

#### Tests First (TDD Red Phase)

- [ ] T039 [US4] Write test US4-AS1 in handlers/todo_handler_test.go: Empty state → helpful message
- [ ] T040 [US4] Write test US4-AS2 in handlers/todo_handler_test.go: 5 todos → all visible
- [ ] T041 [US4] Write test US4-AS3 in handlers/todo_handler_test.go: Persistence across sessions
- [ ] T042 [US4] Write edge case tests: 100+ todos, pagination, sorting by created_at DESC
- [ ] T043 [US4] Run tests with `go test -v ./handlers` → verify all FAIL

#### Implementation (TDD Green Phase)

- [ ] T044 [US4] Implement TodoService.List() in services/todo_service.go with pagination
- [ ] T045 [US4] Implement TodoHandler.List() in handlers/todo_handler.go (GET /api/v1/todos)
- [ ] T046 [US4] Implement TodoHandler.Get() in handlers/todo_handler.go (GET /api/v1/todos/{id})
- [ ] T047 [US4] Register List and Get routes in handlers/routes.go
- [ ] T048 [US4] Run tests with `go test -v ./handlers` → verify all PASS
- [ ] T049 [US4] Run full test suite with `go test -v ./...` → verify no regressions

### Frontend (MVP)

- [ ] T050 [P] [US1] [US4] Create static/index.html with HTML structure (form, list, empty state)
- [ ] T051 [P] [US1] [US4] Add CSS styling in static/index.html (inline or <style> tag)
- [ ] T052 [US1] [US4] Add JavaScript for add todo (fetch POST /api/v1/todos)
- [ ] T053 [US1] [US4] Add JavaScript for list todos (fetch GET /api/v1/todos)
- [ ] T054 [US1] [US4] Add JavaScript for empty state handling
- [ ] T055 [US1] [US4] Add static file serving in cmd/api/main.go
- [ ] T056 [US1] [US4] Manual test: Open http://localhost:8080, add todo, refresh, verify persistence

### MVP Verification

- [ ] T057 Run full test suite with coverage: `go test -v -coverprofile=coverage.out ./...`
- [ ] T058 Verify coverage >80% with `go tool cover -func=coverage.out`
- [ ] T059 Run with race detector: `go test -v -race ./...`
- [ ] T060 Manual end-to-end test of US1 + US4 acceptance scenarios
- [ ] T061 **MVP CHECKPOINT**: Tag release v0.1.0-mvp

---

## Phase 4: User Story 2 (P2) - Mark Complete

**Goal**: Users can track progress by marking todos complete/incomplete

**Duration**: 2-3 hours

**Independent Test**: Add todo, mark complete (strikethrough), mark incomplete (normal), verify visual distinction

**Value Delivered**: Progress tracking functionality

### Tests First (TDD Red Phase)

- [ ] T062 [US2] Write test US2-AS1 in handlers/todo_handler_test.go: Mark complete → completed=true
- [ ] T063 [US2] Write test US2-AS2 in handlers/todo_handler_test.go: Mark incomplete → completed=false
- [ ] T064 [US2] Write test US2-AS3 in handlers/todo_handler_test.go: Mixed states → distinguishable
- [ ] T065 [US2] Write edge case tests: toggle non-existent todo, rapid toggles, filter by completed
- [ ] T066 [US2] Run tests with `go test -v ./handlers` → verify all FAIL

### Implementation (TDD Green Phase)

- [ ] T067 [US2] Implement TodoService.Update() in services/todo_service.go
- [ ] T068 [US2] Implement TodoHandler.Update() in handlers/todo_handler.go (PUT /api/v1/todos/{id})
- [ ] T069 [US2] Register Update route in handlers/routes.go
- [ ] T070 [US2] Run tests with `go test -v ./handlers` → verify all PASS
- [ ] T071 [US2] Run full test suite with `go test -v ./...` → verify no regressions

### Frontend

- [ ] T072 [P] [US2] Add checkbox/button UI for toggle in static/index.html
- [ ] T073 [P] [US2] Add CSS for completed state (strikethrough, gray text) in static/index.html
- [ ] T074 [US2] Add JavaScript for toggle (fetch PUT /api/v1/todos/{id})
- [ ] T075 [US2] Manual test: Toggle completion, verify visual distinction

### Verification

- [ ] T076 Run full test suite with coverage: `go test -v -coverprofile=coverage.out ./...`
- [ ] T077 Verify coverage >80% with `go tool cover -func=coverage.out`
- [ ] T078 Manual end-to-end test of US2 acceptance scenarios

---

## Phase 5: User Story 3 (P3) - Delete Todos

**Goal**: Users can remove irrelevant todos

**Duration**: 1-2 hours

**Independent Test**: Add todo, delete it, verify removed from list

**Value Delivered**: List maintenance functionality

### Tests First (TDD Red Phase)

- [ ] T079 [US3] Write test US3-AS1 in handlers/todo_handler_test.go: Delete todo → removed
- [ ] T080 [US3] Write test US3-AS2 in handlers/todo_handler_test.go: Cancel delete (frontend only, skip backend test)
- [ ] T081 [US3] Write test US3-AS3 in handlers/todo_handler_test.go: Delete completed todo → removed
- [ ] T082 [US3] Write edge case tests: delete non-existent todo, delete twice
- [ ] T083 [US3] Run tests with `go test -v ./handlers` → verify all FAIL

### Implementation (TDD Green Phase)

- [ ] T084 [US3] Implement TodoService.Delete() in services/todo_service.go
- [ ] T085 [US3] Implement TodoHandler.Delete() in handlers/todo_handler.go (DELETE /api/v1/todos/{id})
- [ ] T086 [US3] Register Delete route in handlers/routes.go
- [ ] T087 [US3] Run tests with `go test -v ./handlers` → verify all PASS
- [ ] T088 [US3] Run full test suite with `go test -v ./...` → verify no regressions

### Frontend

- [ ] T089 [P] [US3] Add delete button UI in static/index.html
- [ ] T090 [P] [US3] Add confirmation dialog (optional, for US3-AS2)
- [ ] T091 [US3] Add JavaScript for delete (fetch DELETE /api/v1/todos/{id})
- [ ] T092 [US3] Manual test: Delete todo, verify removed

### Verification

- [ ] T093 Run full test suite with coverage: `go test -v -coverprofile=coverage.out ./...`
- [ ] T094 Verify coverage >80% with `go tool cover -func=coverage.out`
- [ ] T095 Manual end-to-end test of US3 acceptance scenarios

---

## Phase 6: Polish & Cross-Cutting Concerns

**Goal**: Production-ready quality and user experience

**Duration**: 2-3 hours

### Error Handling & Edge Cases

- [ ] T096 [P] Add error handling tests for all error codes in handlers/todo_handler_test.go
- [ ] T097 [P] Test context cancellation in handlers/todo_handler_test.go
- [ ] T098 [P] Test database connection failures in handlers/todo_handler_test.go
- [ ] T099 [P] Add user-friendly error messages in static/index.html JavaScript

### Performance & Optimization

- [ ] T100 [P] Add database indexes (created_at DESC, completed) via migrations
- [ ] T101 [P] Test with 100+ todos, verify <1s load time
- [ ] T102 [P] Add pagination controls in static/index.html (if needed)

### Documentation

- [ ] T103 [P] Update README.md with final setup instructions
- [ ] T104 [P] Add API documentation comments to handlers
- [ ] T105 [P] Create CONTRIBUTING.md with development workflow

### Deployment Preparation

- [ ] T106 [P] Create Dockerfile with multi-stage build
- [ ] T107 [P] Create .dockerignore file
- [ ] T108 [P] Test Docker build: `docker build -t todo-app .`
- [ ] T109 [P] Test Docker run: `docker run -p 8080:8080 todo-app`
- [ ] T110 [P] Create deployment documentation in docs/deployment.md

### Final Verification

- [ ] T111 Run full test suite: `go test -v ./...`
- [ ] T112 Run with race detector: `go test -v -race ./...`
- [ ] T113 Generate coverage report: `go test -v -coverprofile=coverage.out ./...`
- [ ] T114 Verify coverage >80%: `go tool cover -func=coverage.out`
- [ ] T115 Run linter: `golangci-lint run`
- [ ] T116 Manual test all acceptance scenarios (US1-AS1 through US4-AS3)
- [ ] T117 Manual test all edge cases from spec.md
- [ ] T118 Performance test: Add 100 todos, verify <1s operations
- [ ] T119 **FINAL CHECKPOINT**: Tag release v1.0.0

---

## Dependencies & Execution Order

### Critical Path (Must Complete in Order)

```
Phase 1 (Setup)
  ↓
Phase 2 (Foundation) - BLOCKING for all user stories
  ↓
Phase 3 (US1 + US4) - MVP
  ↓
Phase 4 (US2) - Can start after Phase 3
  ↓
Phase 5 (US3) - Can start after Phase 3
  ↓
Phase 6 (Polish) - Requires all phases complete
```

### User Story Dependencies

- **US1 (Add)**: No dependencies (can implement first)
- **US4 (View)**: No dependencies (can implement with US1)
- **US2 (Complete)**: Depends on US1 + US4 (needs todos to mark)
- **US3 (Delete)**: Depends on US1 + US4 (needs todos to delete)

### Parallel Opportunities

**Phase 1**: T003, T004, T005, T006, T007 can be done simultaneously

**Phase 2**: 
- T012, T013, T014 (database layer) can be parallel
- T021, T022 (middleware) can be parallel
- T050, T051 (frontend structure) can start early

**Phase 3**:
- Frontend tasks (T050-T055) can be done while backend tests are being written

**Phase 4-5**:
- US2 and US3 can be implemented in parallel after US1+US4 complete

**Phase 6**:
- Most polish tasks (T096-T110) can be done in parallel

---

## Testing Strategy

### Test Coverage Requirements

- **Minimum**: 80% code coverage for business logic
- **Target**: 90% code coverage overall
- **Focus Areas**:
  - All acceptance scenarios (US#-AS#)
  - All edge cases from spec.md
  - All error conditions
  - Context cancellation
  - Database constraints

### Test Execution

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -v -race ./...

# Run specific test
go test -v -run TestTodoAPI_Create ./handlers
```

### Test Organization

Each test file follows table-driven design:

```go
func TestTodoAPI_Create(t *testing.T) {
    testCases := []struct {
        name     string
        scenario string  // US#-AS# reference
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

---

## Success Criteria

### MVP (Phase 3 Complete)

- ✅ Users can add todos
- ✅ Users can view all todos
- ✅ Todos persist across sessions
- ✅ Empty state shows helpful message
- ✅ Validation prevents empty todos
- ✅ All US1 and US4 acceptance scenarios pass
- ✅ Test coverage >80%

### Complete Feature (All Phases)

- ✅ All user stories implemented (US1-US4)
- ✅ All acceptance scenarios pass (13 scenarios)
- ✅ All edge cases handled
- ✅ Test coverage >80%
- ✅ Performance: <1s for 100 todos
- ✅ Docker image builds successfully
- ✅ Documentation complete

---

## Notes

### TDD Discipline

**CRITICAL**: Follow TDD workflow strictly:
1. Write test FIRST (verify it fails)
2. Implement minimal code to pass
3. Run full test suite
4. Refactor if needed
5. Verify tests still pass

### Constitution Compliance

All tasks follow Go Project Constitution:
- Integration tests with real PostgreSQL (testcontainers)
- Table-driven test design
- ServeHTTP testing via root mux
- Protobuf-only service parameters
- Context-aware operations
- Sentinel errors + HTTP error codes
- Services in public packages

### Incremental Delivery

Each phase delivers working, testable functionality:
- **Phase 3**: MVP (can release to users)
- **Phase 4**: Enhanced (progress tracking)
- **Phase 5**: Complete (full feature set)
- **Phase 6**: Production-ready

### Time Estimates

- **Fast Track** (experienced dev): 1-2 days
- **Normal** (following TDD strictly): 3-4 days
- **Learning** (new to stack): 5-7 days

---

## Quick Reference

### Common Commands

```bash
# Start development
docker-compose up -d postgres
go run cmd/api/main.go

# Run tests
go test -v ./...
go test -v -race ./...
go test -v -coverprofile=coverage.out ./...

# Generate protobuf
go generate ./...

# Build
go build -o bin/todo-app cmd/api/main.go

# Docker
docker build -t todo-app .
docker run -p 8080:8080 todo-app
```

### File Paths Quick Reference

- **Protobuf**: `api/proto/v1/todo.proto`
- **Generated**: `api/gen/v1/*.pb.go`
- **Service**: `services/todo_service.go`
- **Handler**: `handlers/todo_handler.go`
- **Tests**: `handlers/todo_handler_test.go`
- **Model**: `internal/models/todo.go`
- **Frontend**: `static/index.html`
- **Main**: `cmd/api/main.go`

---

**Ready to implement!** Start with Phase 1, follow TDD workflow, and deliver incrementally. 🚀