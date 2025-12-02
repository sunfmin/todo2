package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	todov1 "github.com/yourorg/todo-app/api/gen/v1"
	"github.com/yourorg/todo-app/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// TodoService defines the interface for todo operations
// All methods use protobuf structs (NO primitives)
type TodoService interface {
	Create(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.Todo, error)
	Get(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.Todo, error)
	List(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error)
	Update(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.Todo, error)
	Delete(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error)
}

// todoService implements TodoService
type todoService struct {
	db *gorm.DB
}

// todoServiceBuilder builds a TodoService with optional dependencies
type todoServiceBuilder struct {
	db *gorm.DB
}

// NewTodoService creates a new TodoService builder
// Required parameter: db
func NewTodoService(db *gorm.DB) *todoServiceBuilder {
	return &todoServiceBuilder{db: db}
}

// Build creates the TodoService instance
func (b *todoServiceBuilder) Build() TodoService {
	return &todoService{
		db: b.db,
	}
}

// Create creates a new todo item
// Implements FR-001, FR-006, FR-009, FR-ERR-002, FR-ERR-003
func (s *todoService) Create(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.Todo, error) {
	// FR-006: Trim leading/trailing whitespace while preserving internal whitespace
	// Then validate that trimmed result is not empty
	desc := strings.TrimSpace(req.Description)
	if desc == "" {
		// FR-ERR-003: User-friendly inline validation message
		return nil, fmt.Errorf("create todo: %w", ErrEmptyDescription)
	}
	
	// FR-009: Enforce maximum length of 500 characters (after trimming)
	if len(desc) > 500 {
		return nil, fmt.Errorf("create todo: %w", ErrDescriptionTooLong)
	}

	// Create model
	todo := &models.Todo{
		Description: desc,
		Completed:   false,
	}

	// Save to database (FR-005: Persist todos)
	if err := s.db.WithContext(ctx).Create(todo).Error; err != nil {
		return nil, fmt.Errorf("create todo in database: %w", err)
	}

	return toProto(todo), nil
}

// Get retrieves a single todo by ID
func (s *todoService) Get(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.Todo, error) {
	// Parse UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse todo ID: %w", ErrInvalidInput)
	}

	// Query database
	var todo models.Todo
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("get todo %s: %w", req.Id, ErrTodoNotFound)
		}
		return nil, fmt.Errorf("query todo %s: %w", req.Id, err)
	}

	return toProto(&todo), nil
}

// List retrieves todos with pagination and optional filtering
func (s *todoService) List(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	// Set defaults
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Build query
	query := s.db.WithContext(ctx).Model(&models.Todo{})

	// Apply filter if specified
	if req.Completed != nil {
		query = query.Where("completed = ?", *req.Completed)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count todos: %w", err)
	}

	// Query todos
	var todos []models.Todo
	if err := query.Order("created_at DESC").Limit(int(limit)).Offset(int(offset)).Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	// Convert to protobuf
	pbTodos := make([]*todov1.Todo, len(todos))
	for i, todo := range todos {
		pbTodos[i] = toProto(&todo)
	}

	return &todov1.ListTodosResponse{
		Todos:  pbTodos,
		Total:  int32(total),
		Limit:  limit,
		Offset: offset,
	}, nil
}

// Update updates a todo item
func (s *todoService) Update(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.Todo, error) {
	// Parse UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse todo ID: %w", ErrInvalidInput)
	}

	// Find existing todo
	var todo models.Todo
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("update todo %s: %w", req.Id, ErrTodoNotFound)
		}
		return nil, fmt.Errorf("query todo %s: %w", req.Id, err)
	}

	// Apply updates
	updates := make(map[string]interface{})

	if req.Description != nil {
		// FR-006: Trim leading/trailing whitespace while preserving internal whitespace
		desc := strings.TrimSpace(*req.Description)
		if desc == "" {
			// FR-ERR-003: User-friendly inline validation message
			return nil, fmt.Errorf("update todo: %w", ErrEmptyDescription)
		}
		// FR-009: Enforce maximum length of 500 characters (after trimming)
		if len(desc) > 500 {
			return nil, fmt.Errorf("update todo: %w", ErrDescriptionTooLong)
		}
		updates["description"] = desc
	}

	if req.Completed != nil {
		updates["completed"] = *req.Completed
	}

	// Update in database
	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(&todo).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("update todo %s in database: %w", req.Id, err)
		}
	}

	// Reload to get updated values
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&todo).Error; err != nil {
		return nil, fmt.Errorf("reload todo %s: %w", req.Id, err)
	}

	return toProto(&todo), nil
}

// Delete deletes a todo item
func (s *todoService) Delete(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
	// Parse UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse todo ID: %w", ErrInvalidInput)
	}

	// Delete from database
	result := s.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Todo{})
	if result.Error != nil {
		return nil, fmt.Errorf("delete todo %s: %w", req.Id, result.Error)
	}

	// Check if todo existed
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("delete todo %s: %w", req.Id, ErrTodoNotFound)
	}

	return &todov1.DeleteTodoResponse{}, nil
}

// Helper functions

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