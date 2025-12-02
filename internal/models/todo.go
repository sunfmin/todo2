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