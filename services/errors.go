package services

import "errors"

// Sentinel errors for the service layer
// These are wrapped with context using fmt.Errorf("%w") in service methods
var (
	// ErrTodoNotFound is returned when a todo item is not found
	ErrTodoNotFound = errors.New("todo not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrEmptyDescription is returned when todo description is empty or whitespace-only
	ErrEmptyDescription = errors.New("todo description cannot be empty")
)