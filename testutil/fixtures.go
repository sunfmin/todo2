package testutil

import (
	"github.com/yourorg/todo-app/internal/models"
	"gorm.io/gorm"
)

// CreateTestTodo creates a todo with default values for testing
// Overrides can be provided via the overrides map
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