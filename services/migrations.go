package services

import (
	"github.com/yourorg/todo-app/internal/models"
	"gorm.io/gorm"
)

// AutoMigrate runs database migrations for all models
// This function is exported so external apps can migrate the schema
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Todo{},
	)
}