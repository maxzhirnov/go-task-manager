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
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at"}).
		AddRow(1, "Test Task", "Test Description", "pending", 1, 0, time.Now(), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = ?").
		WithArgs(1).
		WillReturnRows(mockRows)

	handler := NewTaskHandler(db)
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	assert.NoError(t, err)

	req = req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
	rr := httptest.NewRecorder()

	handler.GetTasks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRow := sqlmock.NewRows([]string{"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at"}).
		AddRow(1, "Test Task", "Test Description", "pending", 1, 0, time.Now(), time.Now())
	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnRows(mockRow)

	handler := NewTaskHandler(db)
	req, err := http.NewRequest("GET", "/api/tasks/1", nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.GetTask(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("INSERT INTO tasks").
		WithArgs("Test Task", "Test Description", "pending", 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	handler := NewTaskHandler(db)
	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	taskJSON, err := json.Marshal(task)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(taskJSON))
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))

	rr := httptest.NewRecorder()
	handler.CreateTask(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTaskPositionsInvalidInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := NewTaskHandler(db)
	invalidJSON := []byte(`{"invalid json`)

	req, err := http.NewRequest("PUT", "/api/tasks/positions", bytes.NewBuffer(invalidJSON))
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))

	rr := httptest.NewRecorder()
	handler.UpdateTaskPositions(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTaskPositionsTaskNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	handler := NewTaskHandler(db)
	positions := map[int]int{999: 1}
	positionsJSON, err := json.Marshal(positions)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/tasks/positions", bytes.NewBuffer(positionsJSON))
	assert.NoError(t, err)
	req = req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))

	rr := httptest.NewRecorder()
	handler.UpdateTaskPositions(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
