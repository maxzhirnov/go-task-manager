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
	"github.com/maxzhirnov/go-task-manager/pkg/analytics"
	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	tests := []struct {
		name           string
		setupAuth      func(*http.Request) *http.Request
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
		expectedTasks  []models.Task
	}{
		{
			name: "Successful tasks retrieval",
			setupAuth: func(req *http.Request) *http.Request {
				return req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				createdAt := time.Now()
				updatedAt := time.Now()
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
					}).AddRow(1, "Test Task", "Test Description", "pending", 1, 0, createdAt, updatedAt))
			},
			expectedStatus: http.StatusOK,
			expectedTasks: []models.Task{
				{
					ID:          1,
					Title:       "Test Task",
					Description: "Test Description",
					Status:      "pending",
					UserID:      1,
					Position:    0,
				},
			},
		},
		{
			name: "No tasks found - empty array",
			setupAuth: func(req *http.Request) *http.Request {
				return req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
					}))
			},
			expectedStatus: http.StatusOK,
			expectedTasks:  []models.Task{},
		},
		{
			name: "Missing authentication",
			setupAuth: func(req *http.Request) *http.Request {
				return req // No claims added
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No mock needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
		{
			name: "Database error",
			setupAuth: func(req *http.Request) *http.Request {
				return req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to fetch tasks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Configure mock expectations
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			mockAnalytics := analytics.NewMock("test-key", false)

			// Create handler and request
			handler := NewTaskHandler(db, mockAnalytics)
			req, err := http.NewRequest("GET", "/api/tasks", nil)
			assert.NoError(t, err)

			// Setup authentication if provided
			if tt.setupAuth != nil {
				req = tt.setupAuth(req)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			handler.GetTasks(rr, req)

			// Assert response status
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body
			if tt.expectedError != "" {
				var errorResponse map[string]string
				err = json.NewDecoder(rr.Body).Decode(&errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse["error"])
			} else {
				var tasks []models.Task
				err = json.NewDecoder(rr.Body).Decode(&tasks)
				assert.NoError(t, err)

				// Compare tasks length
				assert.Equal(t, len(tt.expectedTasks), len(tasks))

				// Compare tasks content if any exist
				if len(tt.expectedTasks) > 0 {
					for i, expectedTask := range tt.expectedTasks {
						assert.Equal(t, expectedTask.ID, tasks[i].ID)
						assert.Equal(t, expectedTask.Title, tasks[i].Title)
						assert.Equal(t, expectedTask.Description, tasks[i].Description)
						assert.Equal(t, expectedTask.Status, tasks[i].Status)
						assert.Equal(t, expectedTask.UserID, tasks[i].UserID)
						assert.Equal(t, expectedTask.Position, tasks[i].Position)
					}
				}
			}

			// Verify that all expected mock calls were made
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedError  string
		expectedTask   *models.Task
	}{
		{
			name:   "Successful task retrieval",
			taskID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				createdAt := time.Now()
				updatedAt := time.Now()
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
					}).AddRow(1, "Test Task", "Test Description", "pending", 1, 0, createdAt, updatedAt))
			},
			expectedStatus: http.StatusOK,
			expectedTask: &models.Task{
				ID:          1,
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				UserID:      1,
				Position:    0,
			},
		},
		{
			name:   "Invalid task ID",
			taskID: "invalid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No mock needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid task ID",
		},
		{
			name:   "Task not found",
			taskID: "999",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Task not found",
		},
		{
			name:   "Database error",
			taskID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "sql: connection is already closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Configure mock expectations
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}
			mockAnalytics := analytics.NewMock("test-key", false)
			// Create handler and request
			handler := NewTaskHandler(db, mockAnalytics)
			req, err := http.NewRequest("GET", "/api/tasks/"+tt.taskID, nil)
			assert.NoError(t, err)
			req = mux.SetURLVars(req, map[string]string{"id": tt.taskID})

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			handler.GetTask(rr, req)

			// Assert response status
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body
			if tt.expectedError != "" {
				var errorResponse map[string]string
				err = json.NewDecoder(rr.Body).Decode(&errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse["error"])
			} else if tt.expectedTask != nil {
				var task models.Task
				err = json.NewDecoder(rr.Body).Decode(&task)
				assert.NoError(t, err)
				// Compare relevant fields
				assert.Equal(t, tt.expectedTask.ID, task.ID)
				assert.Equal(t, tt.expectedTask.Title, task.Title)
				assert.Equal(t, tt.expectedTask.Description, task.Description)
				assert.Equal(t, tt.expectedTask.Status, task.Status)
				assert.Equal(t, tt.expectedTask.UserID, task.UserID)
				assert.Equal(t, tt.expectedTask.Position, task.Position)
			}

			// Verify that all expected mock calls were made
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// func TestCreateTask(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		task          interface{} // interface{} to allow invalid JSON
// 		setupAuth     func(*http.Request) *http.Request
// 		mockSetup     func(sqlmock.Sqlmock)
// 		expectedCode  int
// 		expectedError string
// 	}{
// 		{
// 			name: "Successful task creation",
// 			task: models.Task{
// 				Title:       "Test Task",
// 				Description: "Test Description",
// 				Status:      "pending",
// 			},
// 			setupAuth: func(req *http.Request) *http.Request {
// 				return req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
// 			},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				// Expect transaction begin
// 				mock.ExpectBegin()

// 				// Expect update of existing tasks positions
// 				mock.ExpectExec("UPDATE tasks SET position = position \\+ 1").
// 					WithArgs(1).
// 					WillReturnResult(sqlmock.NewResult(0, 2))

// 				// Expect task insertion
// 				mock.ExpectQuery("INSERT INTO tasks \\(title, description, status, user_id, position, created_at, updated_at\\)").
// 					WithArgs(
// 						"Test Task",
// 						"Test Description",
// 						"pending",
// 						1,
// 						0,
// 						sqlmock.AnyArg(),
// 						sqlmock.AnyArg(),
// 					).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 				// Expect transaction commit
// 				mock.ExpectCommit()
// 			},
// 			expectedCode: http.StatusCreated,
// 		},
// 		{
// 			name: "Invalid JSON input",
// 			task: "invalid json",
// 			setupAuth: func(req *http.Request) *http.Request {
// 				return req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
// 			},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				// No database calls expected
// 			},
// 			expectedCode:  http.StatusBadRequest,
// 			expectedError: "Invalid input data",
// 		},
// 		{
// 			name: "Missing title",
// 			task: models.Task{
// 				Description: "Test Description",
// 				Status:      "pending",
// 			},
// 			setupAuth: func(req *http.Request) *http.Request {
// 				return req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
// 			},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				// No database calls expected
// 			},
// 			expectedCode:  http.StatusBadRequest,
// 			expectedError: "Title is required",
// 		},
// 		{
// 			name: "Missing authentication",
// 			task: models.Task{
// 				Title:       "Test Task",
// 				Description: "Test Description",
// 				Status:      "pending",
// 			},
// 			setupAuth: func(req *http.Request) *http.Request {
// 				return req // No claims added
// 			},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				// No database calls expected
// 			},
// 			expectedCode:  http.StatusUnauthorized,
// 			expectedError: "Unauthorized",
// 		},
// 		{
// 			name: "Database error",
// 			task: models.Task{
// 				Title:       "Test Task",
// 				Description: "Test Description",
// 				Status:      "pending",
// 			},
// 			setupAuth: func(req *http.Request) *http.Request {
// 				return req.WithContext(context.WithValue(req.Context(), "claims", &middleware.Claims{UserID: 1}))
// 			},
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectBegin()
// 				mock.ExpectExec("UPDATE tasks SET position = position \\+ 1").
// 					WithArgs(1).
// 					WillReturnError(sql.ErrConnDone)
// 				mock.ExpectRollback()
// 			},
// 			expectedCode:  http.StatusInternalServerError,
// 			expectedError: "Failed to create task",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup
// 			db, mock, err := sqlmock.New()
// 			assert.NoError(t, err)
// 			defer db.Close()

// 			// Configure mock expectations
// 			tt.mockSetup(mock)
// 			mockAnalytics := analytics.NewMock("test-key", false)
// 			// Create handler
// 			handler := NewTaskHandler(db, mockAnalytics)

// 			// Create request
// 			var body []byte
// 			if str, ok := tt.task.(string); ok {
// 				body = []byte(str)
// 			} else {
// 				body, err = json.Marshal(tt.task)
// 				assert.NoError(t, err)
// 			}

// 			req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(body))
// 			assert.NoError(t, err)

// 			// Setup authentication if provided
// 			if tt.setupAuth != nil {
// 				req = tt.setupAuth(req)
// 			}

// 			// Create response recorder
// 			rr := httptest.NewRecorder()

// 			// Execute request
// 			handler.CreateTask(rr, req)

// 			// Assert response
// 			assert.Equal(t, tt.expectedCode, rr.Code)

// 			if tt.expectedError != "" {
// 				var response map[string]string
// 				err = json.NewDecoder(rr.Body).Decode(&response)
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedError, response["error"])
// 			} else {
// 				var createdTask models.Task
// 				err = json.NewDecoder(rr.Body).Decode(&createdTask)
// 				assert.NoError(t, err)
// 				assert.NotZero(t, createdTask.ID)
// 				assert.Equal(t, 0, createdTask.Position)
// 				assert.NotZero(t, createdTask.CreatedAt)
// 				assert.NotZero(t, createdTask.UpdatedAt)
// 			}

// 			// Verify that all expected mock calls were made
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }

func TestUpdateTaskPositionsInvalidInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mockAnalytics := analytics.NewMock("test-key", false)

	handler := NewTaskHandler(db, mockAnalytics)
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

	mockAnalytics := analytics.NewMock("test-key", false)

	handler := NewTaskHandler(db, mockAnalytics)
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

// func TestDeleteTask(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		taskID         string
// 		mockSetup      func(sqlmock.Sqlmock)
// 		expectedStatus int
// 		expectedError  string
// 	}{
// 		{
// 			name:   "Successful status update",
// 			taskID: "1",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec("UPDATE tasks SET status = 'deleted', updated_at = \\$1 WHERE id = \\$2").
// 					WithArgs(sqlmock.AnyArg(), 1).
// 					WillReturnResult(sqlmock.NewResult(0, 1))
// 			},
// 			expectedStatus: http.StatusNoContent,
// 		},
// 		{
// 			name:   "Invalid task ID",
// 			taskID: "invalid",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				// No mock expectations needed for invalid ID
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			expectedError:  "Invalid task ID",
// 		},
// 		{
// 			name:   "Task not found",
// 			taskID: "1",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec("UPDATE tasks SET status = 'deleted', updated_at = \\$1 WHERE id = \\$2").
// 					WithArgs(sqlmock.AnyArg(), 1).
// 					WillReturnResult(sqlmock.NewResult(0, 0))
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedError:  "task not found",
// 		},
// 		{
// 			name:   "Database error",
// 			taskID: "1",
// 			mockSetup: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec("UPDATE tasks SET status = 'deleted', updated_at = \\$1 WHERE id = \\$2").
// 					WithArgs(sqlmock.AnyArg(), 1).
// 					WillReturnError(sql.ErrConnDone)
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedError:  "sql: connection is already closed",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup
// 			db, mock, err := sqlmock.New()
// 			assert.NoError(t, err)
// 			defer db.Close()

// 			// Configure mock expectations
// 			if tt.mockSetup != nil {
// 				tt.mockSetup(mock)
// 			}

// 			mockAnalytics := analytics.NewMock("test-key", false)

// 			// Create handler and request
// 			handler := NewTaskHandler(db, mockAnalytics)
// 			req, err := http.NewRequest("DELETE", "/api/tasks/"+tt.taskID, nil)
// 			assert.NoError(t, err)

// 			// Add URL parameters
// 			req = mux.SetURLVars(req, map[string]string{"id": tt.taskID})

// 			// Create response recorder
// 			rr := httptest.NewRecorder()

// 			// Execute request
// 			handler.DeleteTask(rr, req)

// 			// Assert response status
// 			assert.Equal(t, tt.expectedStatus, rr.Code)

// 			// If we expect an error message, verify it
// 			if tt.expectedError != "" {
// 				var response map[string]string
// 				err = json.NewDecoder(rr.Body).Decode(&response)
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedError, response["error"])
// 			}

// 			// Verify that all expected mock calls were made
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }
