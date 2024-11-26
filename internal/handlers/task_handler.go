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
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Error decoding task: %v", err) // Log the error
		http.Error(w, `{"error": "Invalid input data"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	if task.Status == "" {
		task.Status = "pending"
	}

	// Get the user_id from the JWT claims
	claims, ok := r.Context().Value("claims").(*middleware.Claims)
	if !ok {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	task.UserID = claims.UserID // Associate the task with the user

	if err := task.CreateTask(h.DB); err != nil {
		log.Printf("Error creating task: %v", err) // Log the error
		http.Error(w, `{"error": "Failed to create task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
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
