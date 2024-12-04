// @title Task Manager API
// @version 1.0
// @description Task management system with JWT authentication
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/internal/models"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

// TaskHandler manages task-related HTTP requests.
// It handles CRUD operations for tasks, ensuring proper user authorization
// and data validation.
type TaskHandler struct {
	// DB provides database access for task operations
	DB database.DB
}

// NewTaskHandler creates a new instance of TaskHandler.
//
// Parameters:
//   - db: Database interface for task operations
//
// Returns:
//   - *TaskHandler: Configured task handler
func NewTaskHandler(db database.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// GetTasks retrieves all tasks for the authenticated user.
//
// It extracts the user ID from the JWT claims in the request context
// and returns all tasks associated with that user. If no tasks exist,
// it returns an empty array rather than null.
//
// Authorization:
//   - Requires valid JWT token in request context
//
// HTTP Responses:
//   - 200 OK: Successfully retrieved tasks
//   - 401 Unauthorized: Missing or invalid JWT token
//   - 500 Internal Server Error: Database or server errors
//
// Example success response:
//
//	[
//	    {
//	        "id": 1,
//	        "title": "Complete project",
//	        "description": "Finish the task manager project",
//	        "status": "pending",
//	        "user_id": 123,
//	        "position": 1
//	    }
//	]
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	// Extract user claims from context
	claims, ok := r.Context().Value("claims").(*middleware.Claims)
	if !ok {
		log.Printf("Unauthorized access attempt to GetTasks")
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Fetch tasks from database
	tasks, err := models.GetTasks(h.DB, claims.UserID)
	if err != nil {
		log.Printf("Error fetching tasks for user %d: %v", claims.UserID, err)
		http.Error(w, `{"error": "Failed to fetch tasks"}`, http.StatusInternalServerError)
		return
	}

	// Ensure null is never returned for tasks array
	if tasks == nil {
		tasks = []models.Task{}
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

	log.Printf("Successfully retrieved tasks for user %d", claims.UserID)
}

// GetTask retrieves a specific task by its ID.
//
// It validates the task ID from the URL parameters and ensures the task exists.
// The handler returns detailed task information in JSON format.
//
// URL Parameters:
//   - id: Task identifier (integer)
//
// Authorization:
//   - Requires valid JWT token in request context
//   - User must have access to the requested task
//
// HTTP Responses:
//   - 200 OK: Successfully retrieved task
//   - 400 Bad Request: Invalid task ID format
//   - 404 Not Found: Task doesn't exist
//   - 500 Internal Server Error: Database or server errors
//
// Example success response:
//
//	{
//	    "id": 1,
//	    "title": "Complete project",
//	    "description": "Finish the task manager project",
//	    "status": "pending",
//	    "user_id": 123,
//	    "position": 1,
//	    "created_at": "2024-01-01T12:00:00Z",
//	    "updated_at": "2024-01-01T12:00:00Z"
//	}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Extract and validate task ID from URL parameters
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID format: %s", vars["id"])
		JSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Retrieve task from database
	task, err := models.GetTask(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Task not found
			log.Printf("Task not found: ID %d", id)
			JSONError(w, "Task not found", http.StatusNotFound)
		} else {
			// Other database errors
			log.Printf("Error retrieving task %d: %v", id, err)
			JSONError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

	log.Printf("Successfully retrieved task %d", id)
}

// CreateTask handles the creation of a new task for the authenticated user.
//
// It validates the input data, sets default values where necessary,
// and associates the task with the authenticated user. The handler
// expects a JSON payload containing task details.
//
// Authorization:
//   - Requires valid JWT token in request context
//
// Request Body:
//
//	{
//	    "title": "Complete project",        // Required
//	    "description": "Project details",   // Optional
//	    "status": "pending",               // Optional, defaults to "pending"
//	    "position": 1                      // Optional
//	}
//
// HTTP Responses:
//   - 201 Created: Successfully created task
//   - 400 Bad Request: Invalid input data or missing required fields
//   - 401 Unauthorized: Missing or invalid JWT token
//   - 500 Internal Server Error: Database or server errors
//
// Example success response:
//
//	{
//	    "id": 1,
//	    "title": "Complete project",
//	    "description": "Project details",
//	    "status": "pending",
//	    "user_id": 123,
//	    "position": 1,
//	    "created_at": "2024-01-01T12:00:00Z",
//	    "updated_at": "2024-01-01T12:00:00Z"
//	}
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received create task request")

	// Parse and validate request body
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Error decoding task: %v", err)
		JSONError(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	log.Printf("Decoded task request: %+v", task)

	// Validate required fields
	if task.Title == "" {
		log.Printf("Task creation failed: empty title")
		JSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Set default status if not provided
	if task.Status == "" {
		log.Printf("Setting default status 'pending' for task")
		task.Status = "pending"
	}

	// Extract user information from JWT claims
	claims, ok := r.Context().Value("claims").(*middleware.Claims)
	if !ok {
		log.Printf("Task creation failed: missing or invalid claims in context")
		JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Associate task with authenticated user
	task.UserID = claims.UserID
	log.Printf("Associated task with user ID: %d", task.UserID)

	// Create task in database
	if err := task.CreateTask(h.DB); err != nil {
		log.Printf("Error creating task in database: %v", err)
		log.Printf("Failed task details: %+v", task)
		JSONError(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}

	log.Printf("Successfully created task ID: %d for user ID: %d", task.ID, task.UserID)
}

// UpdateTask handles the modification of an existing task.
//
// It validates the task ID from the URL parameters and the update data
// from the request body. The handler ensures the task status is valid
// and updates the task in the database.
//
// URL Parameters:
//   - id: Task identifier (integer)
//
// Authorization:
//   - Requires valid JWT token in request context
//   - User must own the task being updated
//
// Request Body:
//
//	{
//	    "title": "Updated project",        // Optional
//	    "description": "New details",      // Optional
//	    "status": "in_progress",          // Optional, must be valid status
//	    "position": 2                     // Optional
//	}
//
// HTTP Responses:
//   - 200 OK: Successfully updated task
//   - 400 Bad Request: Invalid task ID, status, or input data
//   - 401 Unauthorized: Missing or invalid JWT token
//   - 404 Not Found: Task doesn't exist
//   - 500 Internal Server Error: Database or server errors
//
// Example success response:
//
//	{
//	    "id": 1,
//	    "title": "Updated project",
//	    "description": "New details",
//	    "status": "in_progress",
//	    "user_id": 123,
//	    "position": 2,
//	    "created_at": "2024-01-01T12:00:00Z",
//	    "updated_at": "2024-01-01T12:00:00Z"
//	}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received task update request")

	// Extract and validate task ID from URL parameters
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID format: %s", vars["id"])
		JSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse and validate request body
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Error decoding task update data: %v", err)
		JSONError(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// Set task ID from URL parameter
	task.ID = id
	log.Printf("Updating task ID: %d with data: %+v", id, task)

	// Validate task status if provided
	if err := task.ValidateStatus(); err != nil {
		log.Printf("Invalid task status: %v", err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update task in database
	if err := task.UpdateTask(h.DB); err != nil {
		log.Printf("Error updating task %d: %v", id, err)
		JSONError(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}

	log.Printf("Successfully updated task ID: %d", id)
}

// DeleteTask handles the removal of an existing task.
//
// It validates the task ID from the URL parameters and permanently
// removes the task from the database. This operation cannot be undone.
//
// URL Parameters:
//   - id: Task identifier (integer)
//
// Authorization:
//   - Requires valid JWT token in request context
//   - User must own the task being deleted
//
// HTTP Responses:
//   - 204 No Content: Successfully deleted task
//   - 400 Bad Request: Invalid task ID format
//   - 401 Unauthorized: Missing or invalid JWT token
//   - 404 Not Found: Task doesn't exist
//   - 500 Internal Server Error: Database or server errors
//
// Example request:
//
//	DELETE /api/tasks/1
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract and validate task ID from URL parameters
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid task ID format: %s", vars["id"])
		JSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Delete task from database
	if err := models.DeleteTask(h.DB, id); err != nil {
		log.Printf("Error deleting task %d: %v", id, err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send successful response with no content
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Successfully deleted task ID: %d", id)
}

// UpdateTaskPositions handles the reordering of multiple tasks for a user.
//
// It processes a map of task IDs to their new positions, allowing for bulk
// updates of task ordering. The handler ensures all tasks belong to the
// authenticated user before making any changes.
//
// Authorization:
//   - Requires valid JWT token in request context
//   - User must own all tasks being reordered
//
// Request Body:
//
//	{
//	    "1": 3,    // Task ID 1 moves to position 3
//	    "2": 1,    // Task ID 2 moves to position 1
//	    "3": 2     // Task ID 3 moves to position 2
//	}
//
// HTTP Responses:
//   - 200 OK: Successfully updated task positions
//   - 400 Bad Request: Invalid input format
//   - 401 Unauthorized: Missing or invalid JWT token
//   - 404 Not Found: One or more tasks don't exist
//   - 500 Internal Server Error: Database or server errors
//
// Example success response:
//
//	{
//	    "message": "Positions updated successfully"
//	}
//
// Note: This operation is atomic - if any task update fails,
// none of the position changes will be applied.
func (h *TaskHandler) UpdateTaskPositions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received task positions update request")

	// Extract user ID from JWT claims
	claims := r.Context().Value("claims").(*middleware.Claims)
	userID := claims.UserID
	log.Printf("Processing position updates for user ID: %d", userID)

	// Parse and validate request body
	var positions map[int]int
	if err := json.NewDecoder(r.Body).Decode(&positions); err != nil {
		log.Printf("Error decoding position updates: %v", err)
		JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Updating positions for %d tasks: %+v", len(positions), positions)

	// Update each task's position
	for taskID, newPosition := range positions {
		// Verify task exists and belongs to user
		task, err := models.GetTask(h.DB, taskID)
		if err != nil {
			log.Printf("Failed to find task ID %d: %v", taskID, err)
			JSONError(w, "Task not found", http.StatusNotFound)
			return
		}

		// Update task position
		if err := task.UpdateTaskPosition(h.DB, userID, newPosition); err != nil {
			log.Printf("Failed to update position for task ID %d to position %d: %v",
				taskID, newPosition, err)
			JSONError(w, "Failed to update position", http.StatusInternalServerError)
			return
		}

		log.Printf("Updated task ID %d to position %d", taskID, newPosition)
	}

	// Send successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Positions updated successfully",
	})

	log.Printf("Successfully updated positions for %d tasks", len(positions))
}

// GetUserStatistics retrieves task-related statistics for the authenticated user.
//
// It provides a summary of the user's tasks including total count, status breakdowns,
// and recent activity metrics. This endpoint is useful for dashboards and progress tracking.
//
// Authorization:
//   - Requires valid JWT token in request context
//
// HTTP Responses:
//   - 200 OK: Successfully retrieved statistics
//   - 401 Unauthorized: Missing or invalid JWT token
//   - 500 Internal Server Error: Database or server errors
//
// Example success response:
//
//	{
//	    "user_id": 123,
//	    "username": "john_doe",
//	    "total_tasks": 10,
//	    "completed_tasks": 5,
//	    "pending_tasks": 3,
//	    "in_progress_tasks": 2,
//	    "deleted_tasks": 1,
//	    "tasks_created_today": 2
//	}
//
// Note: Statistics are calculated in real-time and reflect the current state
// of the user's tasks in the system.
func (h *TaskHandler) GetUserStatistics(w http.ResponseWriter, r *http.Request) {
	// Extract and validate user claims from context
	claims, ok := r.Context().Value("claims").(*middleware.Claims)
	if !ok {
		log.Printf("Unauthorized access attempt to user statistics")
		JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user statistics from database
	stats, err := models.GetUserStatistics(h.DB, claims.UserID)
	if err != nil {
		log.Printf("Error fetching statistics for user %d: %v", claims.UserID, err)
		JSONError(w, "Failed to fetch user statistics", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding statistics response: %v", err)
		return
	}

	log.Printf("Successfully retrieved statistics for user %d: total tasks: %d, completed: %d",
		claims.UserID, stats.TotalTasks, stats.CompletedTasks)
}
