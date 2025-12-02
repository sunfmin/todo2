package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/yourorg/todo-app/api/gen/v1"
	"github.com/yourorg/todo-app/services"
	"github.com/yourorg/todo-app/testutil"
)

// TestMain sets up test environment
func TestMain(m *testing.M) {
	// Run tests
	m.Run()
}

// Test setup helper - returns service, handler, mux, and cleanup
func setupTest(t *testing.T) (services.TodoService, *TodoHandler, http.Handler, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	
	// Create service
	service := services.NewTodoService(db).Build()
	
	// Create handler
	handler := NewTodoHandler(service)
	
	// Setup routes
	mux := SetupRoutes(service)
	
	// Return service, handler, mux, and cleanup function
	return service, handler, mux, func() {
		testutil.TruncateTables(db, "todos")
		cleanup()
	}
}

// Helper to make HTTP requests
func makeRequest(t *testing.T, mux http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error
	
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}
	
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	
	return rr
}

// Helper to decode response
func decodeResponse(t *testing.T, rr *httptest.ResponseRecorder, v interface{}) {
	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
}

// TestTodoAPI_Create tests the Create endpoint (User Story 1)
func TestTodoAPI_Create(t *testing.T) {
	testCases := []struct {
		name        string
		scenario    string
		description string
		wantCode    int
		wantErr     bool
		errContains string
	}{
		{
			name:        "US1-AS1: Add valid todo",
			scenario:    "Given app is open, When user adds 'Buy groceries', Then todo appears",
			description: "Buy groceries",
			wantCode:    http.StatusCreated,
			wantErr:     false,
		},
		{
			name:        "US1-AS2: Add todo to existing list",
			scenario:    "Given app has existing todos, When user adds 'Call dentist', Then new todo appears with existing",
			description: "Call dentist",
			wantCode:    http.StatusCreated,
			wantErr:     false,
		},
		{
			name:        "US1-AS3: Empty todo validation",
			scenario:    "Given user tries to add empty todo, When they submit without text, Then system prevents submission",
			description: "",
			wantCode:    http.StatusBadRequest,
			wantErr:     true,
			errContains: "please enter a task", // Spec requirement (line 77)
		},
		{
			name:        "Edge case: Whitespace-only todo",
			scenario:    "When user tries to add todo with only whitespace, system prevents submission",
			description: "   ",
			wantCode:    http.StatusBadRequest,
			wantErr:     true,
			errContains: "please enter a task", // Spec requirement (line 77)
		},
		{
			name:        "Edge case: Full-width space only",
			scenario:    "When user tries to add todo with only full-width spaces (CJK), system prevents submission",
			description: "　　　", // U+3000 full-width space
			wantCode:    http.StatusBadRequest,
			wantErr:     true,
			errContains: "please enter a task",
		},
		{
			name:        "Edge case: Very long todo (500+ chars)",
			scenario:    "When user enters todo exceeding 500 characters, system handles appropriately",
			description: strings.Repeat("a", 501),
			wantCode:    http.StatusBadRequest,
			wantErr:     true,
			errContains: "must be 500 characters or less", // Spec requirement (line 78)
		},
		{
			name:        "Edge case: Special characters",
			scenario:    "When user enters special characters, system accepts and displays correctly",
			description: "Buy milk & eggs @ store #1 (urgent!)",
			wantCode:    http.StatusCreated,
			wantErr:     false,
		},
		{
			name:        "Edge case: Emojis",
			scenario:    "When user enters emojis, system accepts and displays correctly",
			description: "🛒 Buy groceries 🥕🥛",
			wantCode:    http.StatusCreated,
			wantErr:     false,
		},
		{
			name:        "Edge case: Unicode characters",
			scenario:    "When user enters unicode characters, system accepts and displays correctly",
			description: "学习中文 - Learn Chinese",
			wantCode:    http.StatusCreated,
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, mux, cleanup := setupTest(t)
			defer cleanup()

			// For US1-AS2, create an existing todo first
			if strings.Contains(tc.name, "US1-AS2") {
				existingReq := &pb.CreateTodoRequest{Description: "Existing todo"}
				makeRequest(t, mux, http.MethodPost, "/api/v1/todos", existingReq)
			}

			// Make request
			req := &pb.CreateTodoRequest{Description: tc.description}
			rr := makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)

			// Check status code
			if rr.Code != tc.wantCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.wantCode, rr.Code, rr.Body.String())
			}

			// Check response
			if tc.wantErr {
				var errResp map[string]interface{}
				decodeResponse(t, rr, &errResp)
				
				// Error response format: {code: "...", message: "..."}
				if errMsg, ok := errResp["message"].(string); ok {
					if !strings.Contains(strings.ToLower(errMsg), strings.ToLower(tc.errContains)) {
						t.Errorf("Expected error to contain '%s', got '%s'", tc.errContains, errMsg)
					}
				} else if errCode, ok := errResp["code"].(string); ok {
					// Also check code field
					if !strings.Contains(strings.ToLower(errCode), strings.ToLower(tc.errContains)) {
						t.Errorf("Expected error code to contain '%s', got code='%s', message='%v'", tc.errContains, errCode, errResp["message"])
					}
				} else {
					t.Errorf("Expected error response with 'message' or 'code' field, got: %v", errResp)
				}
			} else {
				var response pb.Todo
				decodeResponse(t, rr, &response)

				// Constitution Principle V: Derive expected from fixtures (NOT response)
				// Only copy truly random fields: UUIDs and timestamps
				expected := &pb.Todo{
					Id:          response.Id,          // Random UUID (copy from response)
					Description: tc.description,       // From request fixture
					Completed:   false,                // Default value for new todos
					CreatedAt:   response.CreatedAt,   // Timestamp (copy from response)
					UpdatedAt:   response.UpdatedAt,   // Timestamp (copy from response)
				}

				// Constitution Principle V: Use protocmp for comparison
				if diff := cmp.Diff(expected, &response, protocmp.Transform()); diff != "" {
					t.Errorf("Todo mismatch (-want +got):\n%s", diff)
				}

				// For US1-AS2, verify both todos exist
				if strings.Contains(tc.name, "US1-AS2") {
					listRr := makeRequest(t, mux, http.MethodGet, "/api/v1/todos", nil)
					var listResp pb.ListTodosResponse
					decodeResponse(t, listRr, &listResp)
					
					if len(listResp.Todos) != 2 {
						t.Errorf("Expected 2 todos, got %d", len(listResp.Todos))
					}
				}
			}
		})
	}
}

// TestTodoAPI_Create_RapidAdditions tests rapid todo additions (edge case)
func TestTodoAPI_Create_RapidAdditions(t *testing.T) {
	_, _, mux, cleanup := setupTest(t)
	defer cleanup()

	// Rapidly add 10 todos
	descriptions := []string{
		"Todo 1", "Todo 2", "Todo 3", "Todo 4", "Todo 5",
		"Todo 6", "Todo 7", "Todo 8", "Todo 9", "Todo 10",
	}

	for _, desc := range descriptions {
		req := &pb.CreateTodoRequest{Description: desc}
		rr := makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)
		
		if rr.Code != http.StatusCreated {
			t.Errorf("Failed to create todo '%s': status %d", desc, rr.Code)
		}
	}

	// Verify all todos were created
	listRr := makeRequest(t, mux, http.MethodGet, "/api/v1/todos", nil)
	var listResp pb.ListTodosResponse
	decodeResponse(t, listRr, &listResp)

	if len(listResp.Todos) != 10 {
		t.Errorf("Expected 10 todos, got %d", len(listResp.Todos))
	}
}

// TestTodoAPI_Create_ContextCancellation tests context cancellation handling
func TestTodoAPI_Create_ContextCancellation(t *testing.T) {
	service, _, _, cleanup := setupTest(t)
	defer cleanup()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Try to create todo with cancelled context
	req := &pb.CreateTodoRequest{Description: "Test todo"}
	_, err := service.Create(ctx, req)

	if err == nil {
		t.Error("Expected error with cancelled context, got nil")
	}
}

// TestTodoAPI_List tests the List endpoint (User Story 4)
func TestTodoAPI_List(t *testing.T) {
	testCases := []struct {
		name        string
		scenario    string
		setupTodos  int
		wantCode    int
		wantCount   int
	}{
		{
			name:       "US4-AS1: Empty state",
			scenario:   "Given no todos exist, When user opens app, Then shows helpful empty message",
			setupTodos: 0,
			wantCode:   http.StatusOK,
			wantCount:  0,
		},
		{
			name:       "US4-AS2: List 5 todos",
			scenario:   "Given 5 todos exist, When user opens app, Then all 5 todos visible",
			setupTodos: 5,
			wantCode:   http.StatusOK,
			wantCount:  5,
		},
		{
			name:       "Edge case: List 100+ todos with pagination",
			scenario:   "When user has 100+ todos, system handles pagination correctly",
			setupTodos: 25,
			wantCode:   http.StatusOK,
			wantCount:  20, // Default limit
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, mux, cleanup := setupTest(t)
			defer cleanup()

			// Setup todos
			for i := 0; i < tc.setupTodos; i++ {
				req := &pb.CreateTodoRequest{
					Description: fmt.Sprintf("Todo %d", i+1),
				}
				makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)
			}

			// List todos
			rr := makeRequest(t, mux, http.MethodGet, "/api/v1/todos", nil)

			if rr.Code != tc.wantCode {
				t.Errorf("Expected status %d, got %d", tc.wantCode, rr.Code)
			}

			var listResp pb.ListTodosResponse
			decodeResponse(t, rr, &listResp)

			if len(listResp.Todos) != tc.wantCount {
				t.Errorf("Expected %d todos, got %d", tc.wantCount, len(listResp.Todos))
			}

			if listResp.Total != int32(tc.setupTodos) {
				t.Errorf("Expected total %d, got %d", tc.setupTodos, listResp.Total)
			}
		})
	}
}

// TestTodoAPI_List_Persistence tests persistence across sessions (US4-AS3)
func TestTodoAPI_List_Persistence(t *testing.T) {
	_, _, mux, cleanup := setupTest(t)
	defer cleanup()

	// Create todos
	descriptions := []string{"Buy milk", "Call dentist", "Finish report"}
	for _, desc := range descriptions {
		req := &pb.CreateTodoRequest{Description: desc}
		makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)
	}

	// Simulate "closing and reopening" by making a new request
	rr := makeRequest(t, mux, http.MethodGet, "/api/v1/todos", nil)

	var listResp pb.ListTodosResponse
	decodeResponse(t, rr, &listResp)

	if len(listResp.Todos) != 3 {
		t.Errorf("Expected 3 persisted todos, got %d", len(listResp.Todos))
	}

	// Verify todos are in correct order (newest first)
	if listResp.Todos[0].Description != "Finish report" {
		t.Errorf("Expected newest todo first, got %s", listResp.Todos[0].Description)
	}
}

// TestTodoAPI_Get tests the Get endpoint
func TestTodoAPI_Get(t *testing.T) {
	testCases := []struct {
		name     string
		scenario string
		setupID  bool
		useID    string
		wantCode int
		wantErr  bool
	}{
		{
			name:     "Get existing todo",
			scenario: "When user requests existing todo, returns todo details",
			setupID:  true,
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name:     "Get non-existent todo",
			scenario: "When user requests non-existent todo, returns 404",
			setupID:  false,
			useID:    "00000000-0000-0000-0000-000000000000",
			wantCode: http.StatusNotFound,
			wantErr:  true,
		},
		{
			name:     "Get with invalid UUID",
			scenario: "When user provides invalid UUID, returns 400",
			setupID:  false,
			useID:    "invalid-uuid",
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, mux, cleanup := setupTest(t)
			defer cleanup()

			var todoID string
			if tc.setupID {
				// Create a todo first
				req := &pb.CreateTodoRequest{Description: "Test todo"}
				createRr := makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)
				var created pb.Todo
				decodeResponse(t, createRr, &created)
				todoID = created.Id
			} else {
				todoID = tc.useID
			}

			// Get the todo
			rr := makeRequest(t, mux, http.MethodGet, fmt.Sprintf("/api/v1/todos/%s", todoID), nil)

			if rr.Code != tc.wantCode {
				t.Errorf("Expected status %d, got %d", tc.wantCode, rr.Code)
			}

			if !tc.wantErr {
				var todo pb.Todo
				decodeResponse(t, rr, &todo)
				if todo.Id != todoID {
					t.Errorf("Expected todo ID %s, got %s", todoID, todo.Id)
				}
			}
		})
	}
}

// TestTodoAPI_Update tests the Update endpoint (User Story 2)
func TestTodoAPI_Update(t *testing.T) {
	testCases := []struct {
		name          string
		scenario      string
		updateReq     func(id string) *pb.UpdateTodoRequest
		wantCode      int
		wantCompleted *bool
		wantErr       bool
	}{
		{
			name:     "US2-AS1: Mark todo complete",
			scenario: "Given incomplete todo, When user marks complete, Then shows strikethrough",
			updateReq: func(id string) *pb.UpdateTodoRequest {
				completed := true
				return &pb.UpdateTodoRequest{Id: id, Completed: &completed}
			},
			wantCode:      http.StatusOK,
			wantCompleted: boolPtr(true),
			wantErr:       false,
		},
		{
			name:     "US2-AS2: Mark todo incomplete",
			scenario: "Given completed todo, When user marks incomplete, Then removes strikethrough",
			updateReq: func(id string) *pb.UpdateTodoRequest {
				completed := false
				return &pb.UpdateTodoRequest{Id: id, Completed: &completed}
			},
			wantCode:      http.StatusOK,
			wantCompleted: boolPtr(false),
			wantErr:       false,
		},
		{
			name:     "Update description",
			scenario: "When user updates todo description, changes are saved",
			updateReq: func(id string) *pb.UpdateTodoRequest {
				desc := "Updated description"
				return &pb.UpdateTodoRequest{Id: id, Description: &desc}
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name:     "Update non-existent todo",
			scenario: "When user tries to update non-existent todo, returns 404",
			updateReq: func(id string) *pb.UpdateTodoRequest {
				completed := true
				return &pb.UpdateTodoRequest{Id: "00000000-0000-0000-0000-000000000000", Completed: &completed}
			},
			wantCode: http.StatusNotFound,
			wantErr:  true,
		},
		{
			name:     "Update with empty description",
			scenario: "When user tries to set empty description, returns validation error with clear message",
			updateReq: func(id string) *pb.UpdateTodoRequest {
				desc := ""
				return &pb.UpdateTodoRequest{Id: id, Description: &desc}
			},
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, mux, cleanup := setupTest(t)
			defer cleanup()

			// Create a todo first
			createReq := &pb.CreateTodoRequest{Description: "Test todo"}
			createRr := makeRequest(t, mux, http.MethodPost, "/api/v1/todos", createReq)
			var created pb.Todo
			decodeResponse(t, createRr, &created)

			// Update the todo
			updateReq := tc.updateReq(created.Id)
			rr := makeRequest(t, mux, http.MethodPut, fmt.Sprintf("/api/v1/todos/%s", updateReq.Id), updateReq)

			if rr.Code != tc.wantCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.wantCode, rr.Code, rr.Body.String())
			}

			if !tc.wantErr {
				var response pb.Todo
				decodeResponse(t, rr, &response)

				// Constitution Principle V: Derive expected from fixtures
				// Build expected based on what was updated
				expected := &pb.Todo{
					Id:        response.Id,        // Random UUID (copy from response)
					CreatedAt: response.CreatedAt, // Timestamp (copy from response)
					UpdatedAt: response.UpdatedAt, // Timestamp (copy from response)
				}

				// Set expected values based on update request
				if updateReq.Description != nil {
					expected.Description = *updateReq.Description
				} else {
					expected.Description = "Test todo" // From CreateTestTodo fixture
				}

				if updateReq.Completed != nil {
					expected.Completed = *updateReq.Completed
				} else {
					expected.Completed = false // Default from fixture
				}

				// Constitution Principle V: Use protocmp for comparison
				if diff := cmp.Diff(expected, &response, protocmp.Transform()); diff != "" {
					t.Errorf("Todo mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// TestTodoAPI_Update_MixedStates tests US2-AS3: Mixed completion states
func TestTodoAPI_Update_MixedStates(t *testing.T) {
	_, _, mux, cleanup := setupTest(t)
	defer cleanup()

	// Create 3 todos
	todos := []string{"Todo 1", "Todo 2", "Todo 3"}
	var todoIDs []string

	for _, desc := range todos {
		req := &pb.CreateTodoRequest{Description: desc}
		rr := makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)
		var created pb.Todo
		decodeResponse(t, rr, &created)
		todoIDs = append(todoIDs, created.Id)
	}

	// Mark first and third as complete
	completed := true
	for _, id := range []string{todoIDs[0], todoIDs[2]} {
		updateReq := &pb.UpdateTodoRequest{Id: id, Completed: &completed}
		makeRequest(t, mux, http.MethodPut, fmt.Sprintf("/api/v1/todos/%s", id), updateReq)
	}

	// List all todos
	listRr := makeRequest(t, mux, http.MethodGet, "/api/v1/todos", nil)
	var listResp pb.ListTodosResponse
	decodeResponse(t, listRr, &listResp)

	// Verify mixed states
	completedCount := 0
	incompleteCount := 0
	for _, todo := range listResp.Todos {
		if todo.Completed {
			completedCount++
		} else {
			incompleteCount++
		}
	}

	if completedCount != 2 {
		t.Errorf("Expected 2 completed todos, got %d", completedCount)
	}
	if incompleteCount != 1 {
		t.Errorf("Expected 1 incomplete todo, got %d", incompleteCount)
	}
}

// TestTodoAPI_Delete tests the Delete endpoint (User Story 3)
func TestTodoAPI_Delete(t *testing.T) {
	testCases := []struct {
		name     string
		scenario string
		setupID  bool
		useID    string
		wantCode int
		wantErr  bool
	}{
		{
			name:     "US3-AS1: Delete todo",
			scenario: "Given todo exists, When user deletes, Then todo removed from list",
			setupID:  true,
			wantCode: http.StatusNoContent,
			wantErr:  false,
		},
		{
			name:     "US3-AS3: Delete completed todo",
			scenario: "Given completed todo, When user deletes, Then todo removed",
			setupID:  true,
			wantCode: http.StatusNoContent,
			wantErr:  false,
		},
		{
			name:     "Delete non-existent todo",
			scenario: "When user tries to delete non-existent todo, returns 404",
			setupID:  false,
			useID:    "00000000-0000-0000-0000-000000000000",
			wantCode: http.StatusNotFound,
			wantErr:  true,
		},
		{
			name:     "Delete with invalid UUID",
			scenario: "When user provides invalid UUID, returns 400",
			setupID:  false,
			useID:    "invalid-uuid",
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, mux, cleanup := setupTest(t)
			defer cleanup()

			var todoID string
			if tc.setupID {
				// Create a todo first
				req := &pb.CreateTodoRequest{Description: "Test todo"}
				createRr := makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)
				var created pb.Todo
				decodeResponse(t, createRr, &created)
				todoID = created.Id

				// For US3-AS3, mark it complete first
				if strings.Contains(tc.name, "US3-AS3") {
					completed := true
					updateReq := &pb.UpdateTodoRequest{Id: todoID, Completed: &completed}
					makeRequest(t, mux, http.MethodPut, fmt.Sprintf("/api/v1/todos/%s", todoID), updateReq)
				}
			} else {
				todoID = tc.useID
			}

			// Delete the todo
			rr := makeRequest(t, mux, http.MethodDelete, fmt.Sprintf("/api/v1/todos/%s", todoID), nil)

			if rr.Code != tc.wantCode {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.wantCode, rr.Code, rr.Body.String())
			}

			// Verify todo was deleted (for successful cases)
			if !tc.wantErr && tc.setupID {
				listRr := makeRequest(t, mux, http.MethodGet, "/api/v1/todos", nil)
				var listResp pb.ListTodosResponse
				decodeResponse(t, listRr, &listResp)

				for _, todo := range listResp.Todos {
					if todo.Id == todoID {
						t.Error("Todo should have been deleted but still exists")
					}
				}
			}
		})
	}
}

// TestTodoAPI_Delete_Twice tests deleting the same todo twice
func TestTodoAPI_Delete_Twice(t *testing.T) {
	_, _, mux, cleanup := setupTest(t)
	defer cleanup()

	// Create a todo
	req := &pb.CreateTodoRequest{Description: "Test todo"}
	createRr := makeRequest(t, mux, http.MethodPost, "/api/v1/todos", req)
	var created pb.Todo
	decodeResponse(t, createRr, &created)

	// Delete once (should succeed)
	rr1 := makeRequest(t, mux, http.MethodDelete, fmt.Sprintf("/api/v1/todos/%s", created.Id), nil)
	if rr1.Code != http.StatusNoContent {
		t.Errorf("First delete: expected status %d, got %d", http.StatusNoContent, rr1.Code)
	}

	// Delete again (should fail with 404)
	rr2 := makeRequest(t, mux, http.MethodDelete, fmt.Sprintf("/api/v1/todos/%s", created.Id), nil)
	if rr2.Code != http.StatusNotFound {
		t.Errorf("Second delete: expected status %d, got %d", http.StatusNotFound, rr2.Code)
	}
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
}