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
	// Spec requirement (line 77): "Please enter a task"
	ErrEmptyDescription = errors.New("please enter a task")

	// ErrDescriptionTooLong is returned when description exceeds 500 characters
	ErrDescriptionTooLong = errors.New("todo description must be 500 characters or less")
)