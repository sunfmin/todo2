# Data Model: Simple Todo App

**Feature**: 001-todo-app  
**Date**: 2025-12-02  
**Purpose**: Define entities, relationships, and validation rules

## Overview

This document defines the data model for the todo application, including entity structures, database schema, validation rules, and state transitions.

## Entities

### Todo Item

**Description**: Represents a single task that a user needs to complete.

**Properties**:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | Primary Key, Auto-generated | Unique identifier |
| `description` | String | NOT NULL, Max 500 chars | Task description |
| `completed` | Boolean | NOT NULL, Default false | Completion status |
| `created_at` | Timestamp | NOT NULL, Auto-generated | Creation timestamp |
| `updated_at` | Timestamp | NOT NULL, Auto-updated | Last modification timestamp |

**Validation Rules** (from FR-006, Edge Cases):
- `description` MUST NOT be empty or whitespace-only
- `description` MUST NOT exceed 500 characters
- `description` MAY contain special characters and emojis
- `completed` MUST be boolean (true/false)

**Indexes**:
- Primary key on `id` (automatic)
- Index on `created_at` for sorting (newest first)
- Index on `completed` for filtering (optional optimization)

**State Transitions**:
```
[New] --create--> [Active (completed=false)]
[Active] --toggle--> [Completed (completed=true)]
[Completed] --toggle--> [Active (completed=false)]
[Active|Completed] --delete--> [Deleted (removed from DB)]
```

## Database Schema

### PostgreSQL Table Definition

```sql
CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description VARCHAR(500) NOT NULL CHECK (LENGTH(TRIM(description)) > 0),
    completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for sorting by creation time (newest first)
CREATE INDEX idx_todos_created_at ON todos(created_at DESC);

-- Optional: Index for filtering by completion status
CREATE INDEX idx_todos_completed ON todos(completed);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_todos_updated_at BEFORE UPDATE ON todos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### GORM Model (Internal)

**Location**: `internal/models/todo.go`

```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Todo represents a task item in the database
// This is an INTERNAL model - services return protobuf types
type Todo struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Description string    `gorm:"type:varchar(500);not null;check:length(trim(description)) > 0"`
    Completed   bool      `gorm:"not null;default:false"`
    CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
    UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (Todo) TableName() string {
    return "todos"
}

// BeforeCreate hook to ensure ID is set
func (t *Todo) BeforeCreate(tx *gorm.DB) error {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    return nil
}
```

## Protobuf Definitions (Public API)

**Location**: `api/proto/v1/todo.proto`

These are the PUBLIC types that services return and handlers use.

```protobuf
syntax = "proto3";

package todo.v1;

option go_package = "github.com/yourorg/todo-app/api/gen/v1;todov1";

import "google/protobuf/timestamp.proto";
import "validate/validate.proto";

// Todo represents a task item
message Todo {
    string id = 1;
    string description = 2;
    bool completed = 3;
    google.protobuf.Timestamp created_at = 4;
    google.protobuf.Timestamp updated_at = 5;
}

// CreateTodoRequest for creating a new todo
message CreateTodoRequest {
    string description = 1 [
        (validate.rules).string = {
            min_len: 1,
            max_len: 500,
            pattern: "^\\S.*$"  // Must not start with whitespace
        }
    ];
}

// GetTodoRequest for retrieving a single todo
message GetTodoRequest {
    string id = 1 [(validate.rules).string.uuid = true];
}

// UpdateTodoRequest for updating a todo
message UpdateTodoRequest {
    string id = 1 [(validate.rules).string.uuid = true];
    optional string description = 2 [
        (validate.rules).string = {
            min_len: 1,
            max_len: 500,
            pattern: "^\\S.*$"
        }
    ];
    optional bool completed = 3;
}

// DeleteTodoRequest for deleting a todo
message DeleteTodoRequest {
    string id = 1 [(validate.rules).string.uuid = true];
}

// ListTodosRequest for listing todos with pagination
message ListTodosRequest {
    int32 limit = 1 [
        (validate.rules).int32 = {
            gte: 1,
            lte: 100
        }
    ];
    int32 offset = 2 [(validate.rules).int32.gte = 0];
    optional bool completed = 3;  // Filter by completion status
}

// ListTodosResponse contains paginated todos
message ListTodosResponse {
    repeated Todo todos = 1;
    int32 total = 2;
    int32 limit = 3;
    int32 offset = 4;
}

// Empty response for delete operation
message DeleteTodoResponse {}
```

## Model Conversion

**Location**: `services/todo_service.go` (internal helper functions)

```go
// toProto converts internal GORM model to public protobuf type
func toProto(t *models.Todo) *todov1.Todo {
    return &todov1.Todo{
        Id:          t.ID.String(),
        Description: t.Description,
        Completed:   t.Completed,
        CreatedAt:   timestamppb.New(t.CreatedAt),
        UpdatedAt:   timestamppb.New(t.UpdatedAt),
    }
}

// fromCreateRequest converts protobuf request to internal model
func fromCreateRequest(req *todov1.CreateTodoRequest) *models.Todo {
    return &models.Todo{
        Description: strings.TrimSpace(req.Description),
        Completed:   false,
    }
}

// applyUpdate applies protobuf update request to internal model
func applyUpdate(t *models.Todo, req *todov1.UpdateTodoRequest) {
    if req.Description != nil {
        t.Description = strings.TrimSpace(*req.Description)
    }
    if req.Completed != nil {
        t.Completed = *req.Completed
    }
}
```

## Validation Strategy

### Input Validation (Protobuf)

Validation happens at the protobuf level using `protoc-gen-validate`:
- Description: 1-500 chars, no leading whitespace
- ID: Valid UUID format
- Limit: 1-100 (pagination)
- Offset: >= 0

### Business Logic Validation (Service Layer)

Additional validation in service methods:
- Check if todo exists before update/delete
- Trim whitespace from description
- Validate UUID format

### Database Constraints

PostgreSQL enforces:
- NOT NULL constraints
- CHECK constraint on description length
- UUID format for ID
- Timestamp defaults

## Query Patterns

### Common Queries

**List All Todos (Paginated)**:
```go
db.WithContext(ctx).
    Order("created_at DESC").
    Limit(limit).
    Offset(offset).
    Find(&todos)
```

**List by Completion Status**:
```go
db.WithContext(ctx).
    Where("completed = ?", completed).
    Order("created_at DESC").
    Limit(limit).
    Offset(offset).
    Find(&todos)
```

**Get Single Todo**:
```go
db.WithContext(ctx).
    Where("id = ?", id).
    First(&todo)
```

**Create Todo**:
```go
db.WithContext(ctx).
    Create(&todo)
```

**Update Todo**:
```go
db.WithContext(ctx).
    Model(&todo).
    Updates(map[string]interface{}{
        "description": newDescription,
        "completed": newCompleted,
    })
```

**Delete Todo**:
```go
db.WithContext(ctx).
    Delete(&todo)
```

**Count Total Todos**:
```go
db.WithContext(ctx).
    Model(&models.Todo{}).
    Count(&total)
```

## Performance Considerations

### Indexing Strategy

1. **Primary Key (id)**: Automatic B-tree index for fast lookups
2. **created_at DESC**: Index for sorting newest first (most common query)
3. **completed**: Optional index if filtering by status is frequent

### Query Optimization

- Use pagination (limit/offset) to prevent loading all todos
- Index on `created_at` supports ORDER BY without table scan
- Avoid SELECT * - only fetch needed columns (GORM handles this)

### Scalability

For 10,000 todos per user:
- Table size: ~1MB (assuming 100 bytes per row)
- Index size: ~200KB (id + created_at indexes)
- Query time: <10ms with proper indexes
- No caching needed at this scale

## Migration Strategy

### Initial Migration

Use GORM AutoMigrate for development:

```go
// services/migrations.go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(&models.Todo{})
}
```

### Production Migrations

For production, consider using migration tools:
- golang-migrate/migrate
- goose
- Or GORM AutoMigrate with version tracking

## Testing Data

### Test Fixtures

**Location**: `testutil/fixtures.go`

```go
// CreateTestTodo creates a todo with default values for testing
func CreateTestTodo(db *gorm.DB, overrides map[string]interface{}) *models.Todo {
    todo := &models.Todo{
        Description: "Test todo item",
        Completed:   false,
    }
    
    // Apply overrides
    if desc, ok := overrides["description"].(string); ok {
        todo.Description = desc
    }
    if completed, ok := overrides["completed"].(bool); ok {
        todo.Completed = completed
    }
    
    db.Create(todo)
    return todo
}
```

### Test Scenarios

From spec.md acceptance scenarios:
- US1-AS1: Create todo "Buy groceries"
- US1-AS2: Create multiple todos
- US1-AS3: Attempt empty todo (validation error)
- US2-AS1: Toggle completion status
- US2-AS2: Toggle back to incomplete
- US2-AS3: Mixed completion states
- US3-AS1: Delete todo
- US3-AS2: Cancel delete (not applicable for API)
- US3-AS3: Delete completed todo
- US4-AS1: Empty state
- US4-AS2: List 5 todos
- US4-AS3: Persistence across sessions

## Relationships

**Current**: None (single entity)

**Future Considerations** (out of scope for MVP):
- User entity (for multi-user support)
- Category/Tag entities (for organization)
- Attachment entity (for files)
- Comment entity (for notes)

## Summary

- **Single entity**: Todo with 5 fields
- **Simple schema**: No foreign keys or complex relationships
- **Validation**: Three layers (protobuf, service, database)
- **Indexing**: Optimized for common queries
- **Scalability**: Handles 10,000 todos efficiently
- **Testing**: Fixtures support all acceptance scenarios