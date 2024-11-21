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

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		JSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = id
	if err := task.UpdateTask(h.DB); err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

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
