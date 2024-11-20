package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/maxzhirnov/go-task-manager/internal/handlers"
	"github.com/stretchr/testify/assert"
)

var mockDB *sql.DB
var mock sqlmock.Sqlmock

// Mock the InitDB function to return sqlmock
func mockInitDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

// Create a custom setupRouter for testing that uses the mock database
func setupMockRouter() *mux.Router {
	// Mock the database connection
	mockDB, mock, _ = mockInitDB() // Mocked DB

	// Create router
	r := mux.NewRouter()

	// Create handler with the mock DB
	taskHandler := handlers.NewTaskHandler(mockDB)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// Serve a default route, e.g., a homepage or a health check
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Task Manager API"))
	})

	return r
}

func TestMainApp(t *testing.T) {
	// Use the mock router
	r := setupMockRouter()

	// Mock the expected SQL query and result for the test
	now := time.Now()
	mock.ExpectQuery("SELECT id, title, description, status, created_at, updated_at FROM tasks").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(1, "Test Task", "Test Description", "pending", now, now))

	// Create a test request to hit an existing route (e.g., `/api/tasks`)
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the test request using the mock router
	r.ServeHTTP(rr, req)

	// Assert the status code (expecting 200 because the route exists)
	assert.Equal(t, http.StatusOK, rr.Code)
}
