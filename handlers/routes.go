package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourorg/todo-app/services"
)

// SetupRoutes creates the HTTP router with all routes registered
// CRITICAL: Production and tests MUST use the SAME routing configuration
func SetupRoutes(service services.TodoService) http.Handler {
	mux := http.NewServeMux()
	handler := NewTodoHandler(service)

	// API routes
	mux.HandleFunc("POST /api/v1/todos", handler.Create)
	mux.HandleFunc("GET /api/v1/todos", handler.List)
	mux.HandleFunc("GET /api/v1/todos/{id}", handler.Get)
	mux.HandleFunc("PUT /api/v1/todos/{id}", handler.Update)
	mux.HandleFunc("DELETE /api/v1/todos/{id}", handler.Delete)

	// Health check
	mux.HandleFunc("GET /health", healthCheck)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /", fs)

	return mux
}

// healthCheck handles the health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}