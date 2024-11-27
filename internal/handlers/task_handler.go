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

type TaskHandler struct {
	DB database.DB
}

func NewTaskHandler(db database.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// @Summary Get all tasks
// @Description Get all tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Task
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /tasks [get]
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	// Get the user_id from the JWT claims
	claims, ok := r.Context().Value("claims").(*middleware.Claims)
	if !ok {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	tasks, err := models.GetTasks(h.DB, claims.UserID) // Fetch tasks for the specific user
	if err != nil {
		log.Printf("Error fetching tasks: %v", err) // Log the error
		http.Error(w, `{"error": "Failed to fetch tasks"}`, http.StatusInternalServerError)
		return
	}

	// If no tasks exist, return an empty array
	if tasks == nil {
		tasks = []models.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// @Summary Get a task
// @Description Get a specific task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Security BearerAuth
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse "Invalid task ID"
// @Failure 404 {object} models.ErrorResponse "Task not found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		JSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := models.GetTask(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			JSONError(w, "Task not found", http.StatusNotFound)
		} else {
			JSONError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// @Summary Create a task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task object"
// @Security BearerAuth
// @Success 201 {object} models.Task
// @Failure 400 {object} models.ErrorResponse "Invalid input data"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received create task request")

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Error decoding task: %v", err)
		http.Error(w, `{"error": "Invalid input data"}`, http.StatusBadRequest)
		return
	}
	log.Printf("Decoded task: %+v", task)

	if task.Title == "" {
		log.Printf("Task creation failed: empty title")
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	if task.Status == "" {
		log.Printf("Setting default status 'pending' for task")
		task.Status = "pending"
	}

	// Get the user_id from the JWT claims
	claims, ok := r.Context().Value("claims").(*middleware.Claims)
	if !ok {
		log.Printf("Task creation failed: missing or invalid claims in context")
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	task.UserID = claims.UserID
	log.Printf("Associated task with user ID: %d", task.UserID)

	if err := task.CreateTask(h.DB); err != nil {
		log.Printf("Error creating task in database: %v", err)
		log.Printf("Failed task details: %+v", task)
		http.Error(w, `{"error": "Failed to create task"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created task with ID: %d", task.ID)
	log.Printf("Final task details: %+v", task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task object"
// @Security BearerAuth
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse "Invalid input data"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received task update request")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		JSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		JSONError(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	task.ID = id // Assign the ID from the request URL

	if err := task.ValidateStatus(); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := task.UpdateTask(h.DB); err != nil {
		JSONError(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse "Invalid task ID"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		JSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteTask(h.DB, id); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update task positions
// @Description Update the positions of multiple tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param positions body map[int]int true "Map of task IDs to new positions"
// @Success 200 {object} map[string]string "Positions updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 404 {object} models.ErrorResponse "Task not found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /tasks/positions [put]
func (h *TaskHandler) UpdateTaskPositions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received positions update request")

	claims := r.Context().Value("claims").(*middleware.Claims)
	userID := claims.UserID

	var positions map[int]int
	if err := json.NewDecoder(r.Body).Decode(&positions); err != nil {
		log.Printf("Error decoding positions: %v", err)
		JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Parsed positions map: %+v", positions)

	for taskID, newPosition := range positions {
		// Get the task first
		task, err := models.GetTask(h.DB, taskID)
		if err != nil {
			JSONError(w, "Task not found", http.StatusNotFound)
			return
		}

		// Update the task's position
		if err := task.UpdateTaskPosition(h.DB, userID, newPosition); err != nil {
			JSONError(w, "Failed to update position", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Positions updated successfully"})
}
