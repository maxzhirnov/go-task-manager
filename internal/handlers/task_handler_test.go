package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mockRows := sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
		AddRow(1, "Test Task", "Test Description", "pending", time.Now(), time.Now())
	mock.ExpectQuery("SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE user_id = ?").
		WithArgs(1). // Mocking for user_id = 1
		WillReturnRows(mockRows)

	// Initialize the handler
	handler := NewTaskHandler(db)

	// Create a test request
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	assert.NoError(t, err)

	// Add a mock user_id to the context to simulate JWT claims
	req = req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{
		UserID: 1,
	}))

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.GetTasks(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body
	var tasks []models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Test Task", tasks[0].Title)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTask(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mockRow := sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
		AddRow(1, "Test Task", "Test Description", "pending", time.Now(), time.Now())
	mock.ExpectQuery("SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnRows(mockRow)

	// Initialize the handler
	handler := NewTaskHandler(db)

	// Create a test request with URL parameters
	req, err := http.NewRequest("GET", "/api/tasks/1", nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Add a mock user_id to the context to simulate JWT claims
	req = req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{
		UserID: 1,
	}))

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.GetTask(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body
	var task models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &task)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task", task.Title)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTask(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mock.ExpectQuery("INSERT INTO tasks").
		WithArgs("Test Task", "Test Description", "pending", 1, sqlmock.AnyArg(), sqlmock.AnyArg()). // Include user_id = 1
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Initialize the handler
	handler := NewTaskHandler(db)

	// Create a test task
	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	// Marshal the task into JSON
	taskJSON, err := json.Marshal(task)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(taskJSON))
	assert.NoError(t, err)

	// Add a mock user_id to the context to simulate JWT claims
	req = req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{
		UserID: 1,
	}))

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.CreateTask(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Assert the response body
	var createdTask models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.Equal(t, 1, createdTask.ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTask(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mock.ExpectExec("UPDATE tasks").
		WithArgs("Updated Task", "Updated Description", "completed", sqlmock.AnyArg(), 1). // Include updated_at timestamp
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Initialize the handler
	handler := NewTaskHandler(db)

	// Create a test task
	task := models.Task{
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "completed",
	}

	// Marshal the task into JSON
	taskJSON, err := json.Marshal(task)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest("PUT", "/api/tasks/1", bytes.NewBuffer(taskJSON))
	assert.NoError(t, err)

	// Add URL parameters to request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.UpdateTask(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTask(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mock.ExpectExec("DELETE FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Initialize the handler
	handler := NewTaskHandler(db)

	// Create a test request
	req, err := http.NewRequest("DELETE", "/api/tasks/1", nil)
	assert.NoError(t, err)

	// Add URL parameters to request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.DeleteTask(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTaskNotFound(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query to return no rows
	mock.ExpectQuery("SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?").
		WithArgs(999). // Use a task ID that doesn't exist
		WillReturnError(sql.ErrNoRows)

	// Initialize the handler
	handler := NewTaskHandler(db)

	// Create a test request with URL parameters
	req, err := http.NewRequest("GET", "/api/tasks/999", nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.GetTask(rr, req)

	// Assert the status code (expecting 404 for task not found)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
