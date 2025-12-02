package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	todov1 "github.com/yourorg/todo-app/api/gen/v1"
	"github.com/yourorg/todo-app/services"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	service services.TodoService
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(service services.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// Create handles POST /api/v1/todos
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req todov1.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, Errors.InvalidRequest)
		return
	}

	todo, err := h.service.Create(r.Context(), &req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// List handles GET /api/v1/todos
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	
	req := &todov1.ListTodosRequest{
		Limit:  20, // default
		Offset: 0,  // default
	}

	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		var limit int32
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil {
			req.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := query.Get("offset"); offsetStr != "" {
		var offset int32
		if _, err := fmt.Sscanf(offsetStr, "%d", &offset); err == nil {
			req.Offset = offset
		}
	}

	// Parse completed filter
	if completedStr := query.Get("completed"); completedStr != "" {
		if completedStr == "true" {
			completed := true
			req.Completed = &completed
		} else if completedStr == "false" {
			completed := false
			req.Completed = &completed
		}
	}

	response, err := h.service.List(r.Context(), req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get handles GET /api/v1/todos/{id}
func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		RespondWithError(w, Errors.InvalidRequest)
		return
	}

	req := &todov1.GetTodoRequest{Id: id}
	todo, err := h.service.Get(r.Context(), req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// Update handles PUT /api/v1/todos/{id}
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		RespondWithError(w, Errors.InvalidRequest)
		return
	}

	var req todov1.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, Errors.InvalidRequest)
		return
	}

	req.Id = id

	todo, err := h.service.Update(r.Context(), &req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// Delete handles DELETE /api/v1/todos/{id}
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		RespondWithError(w, Errors.InvalidRequest)
		return
	}

	req := &todov1.DeleteTodoRequest{Id: id}
	_, err := h.service.Delete(r.Context(), req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}