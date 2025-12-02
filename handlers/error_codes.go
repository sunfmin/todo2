package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/yourorg/todo-app/services"
)

// ErrorCode represents an HTTP error response
type ErrorCode struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	ServiceErr error  `json:"-"` // Maps to service sentinel error
}

// Errors is a singleton containing all error codes
// Messages aligned with spec requirements (FR-ERR-003)
var Errors = struct {
	InvalidRequest      ErrorCode
	TodoNotFound        ErrorCode
	EmptyDescription    ErrorCode
	DescriptionTooLong  ErrorCode
	InternalError       ErrorCode
}{
	InvalidRequest: ErrorCode{
		Code:       "INVALID_REQUEST",
		Message:    "Invalid request data",
		HTTPStatus: http.StatusBadRequest,
		ServiceErr: services.ErrInvalidInput,
	},
	TodoNotFound: ErrorCode{
		Code:       "TODO_NOT_FOUND",
		Message:    "Todo not found",
		HTTPStatus: http.StatusNotFound,
		ServiceErr: services.ErrTodoNotFound,
	},
	EmptyDescription: ErrorCode{
		Code:       "EMPTY_DESCRIPTION",
		Message:    "please enter a task", // Spec requirement (line 77) - lowercase to match service error
		HTTPStatus: http.StatusBadRequest,
		ServiceErr: services.ErrEmptyDescription,
	},
	DescriptionTooLong: ErrorCode{
		Code:       "DESCRIPTION_TOO_LONG",
		Message:    "todo description must be 500 characters or less", // Spec requirement (line 78)
		HTTPStatus: http.StatusBadRequest,
		ServiceErr: services.ErrDescriptionTooLong,
	},
	InternalError: ErrorCode{
		Code:       "INTERNAL_ERROR",
		Message:    "An unexpected error occurred",
		HTTPStatus: http.StatusInternalServerError,
		ServiceErr: nil,
	},
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, errCode ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode.HTTPStatus)
	json.NewEncoder(w).Encode(errCode)
}

// HandleServiceError automatically maps service errors to HTTP responses
// Implements FR-ERR-001: Provide clear feedback when operations fail
func HandleServiceError(w http.ResponseWriter, err error) {
	// Check context errors first
	if errors.Is(err, context.Canceled) {
		w.WriteHeader(499) // Client Closed Request
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	// Check service error mapping
	// Order matters: check more specific errors first
	allErrors := []ErrorCode{
		Errors.EmptyDescription,
		Errors.DescriptionTooLong,
		Errors.TodoNotFound,
		Errors.InvalidRequest,
	}

	for _, errCode := range allErrors {
		if errCode.ServiceErr != nil && errors.Is(err, errCode.ServiceErr) {
			RespondWithError(w, errCode)
			return
		}
	}

	// Default to internal error
	RespondWithError(w, Errors.InternalError)
}