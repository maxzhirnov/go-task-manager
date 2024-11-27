package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		mockSetup     func(sqlmock.Sqlmock)
		expectedTasks []Task
		expectedError error
		errorCheck    func(error) bool // Add a custom error check function
	}{
		{
			name:   "Successfully get multiple tasks",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
				}).
					AddRow(1, "Task 1", "Description 1", "pending", 1, 0, time.Now(), time.Now()).
					AddRow(2, "Task 2", "Description 2", "in_progress", 1, 1, time.Now(), time.Now())

				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = \\$1 ORDER BY position ASC").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedTasks: []Task{
				{
					ID:          1,
					Title:       "Task 1",
					Description: "Description 1",
					Status:      "pending",
					UserID:      1,
					Position:    0,
				},
				{
					ID:          2,
					Title:       "Task 2",
					Description: "Description 2",
					Status:      "in_progress",
					UserID:      1,
					Position:    1,
				},
			},
			expectedError: nil,
		},
		{
			name:   "No tasks found",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
				})

				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = \\$1 ORDER BY position ASC").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedTasks: []Task{},
			expectedError: nil,
		},
		{
			name:   "Database error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = \\$1 ORDER BY position ASC").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedTasks: nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:   "Row scan error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
				}).
					AddRow("invalid", "Task 1", "Description 1", "pending", 1, 0, time.Now(), time.Now()) // Invalid ID type

				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE user_id = \\$1 ORDER BY position ASC").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedTasks: nil,
			expectedError: fmt.Errorf("scan error"), // Just need any error
			errorCheck: func(err error) bool {
				return err != nil // Just check that there is an error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Configure mock expectations
			tt.mockSetup(mock)

			// Execute function
			tasks, err := GetTasks(db, tt.userID)

			// Assert error
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.errorCheck != nil {
					assert.True(t, tt.errorCheck(err))
				} else if tt.expectedError == sql.ErrConnDone {
					assert.Equal(t, tt.expectedError, err)
				}
				return
			}
			assert.NoError(t, err)

			// Assert tasks
			assert.Equal(t, len(tt.expectedTasks), len(tasks))
			for i, expectedTask := range tt.expectedTasks {
				assert.Equal(t, expectedTask.ID, tasks[i].ID)
				assert.Equal(t, expectedTask.Title, tasks[i].Title)
				assert.Equal(t, expectedTask.Description, tasks[i].Description)
				assert.Equal(t, expectedTask.Status, tasks[i].Status)
				assert.Equal(t, expectedTask.UserID, tasks[i].UserID)
				assert.Equal(t, expectedTask.Position, tasks[i].Position)
				assert.NotZero(t, tasks[i].CreatedAt)
				assert.NotZero(t, tasks[i].UpdatedAt)
			}

			// Verify that all expected mock calls were made
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Test Task struct JSON marshaling/unmarshaling
func TestTask_JSON(t *testing.T) {
	now := time.Now()
	task := Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      StatusPending,
		UserID:      1,
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test marshaling
	jsonData, err := json.Marshal(task)
	assert.NoError(t, err)

	// Test unmarshaling
	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	assert.NoError(t, err)

	// Compare fields
	assert.Equal(t, task.ID, unmarshaledTask.ID)
	assert.Equal(t, task.Title, unmarshaledTask.Title)
	assert.Equal(t, task.Description, unmarshaledTask.Description)
	assert.Equal(t, task.Status, unmarshaledTask.Status)
	assert.Equal(t, task.UserID, unmarshaledTask.UserID)
	assert.Equal(t, task.Position, unmarshaledTask.Position)
	assert.Equal(t, task.CreatedAt.Unix(), unmarshaledTask.CreatedAt.Unix())
	assert.Equal(t, task.UpdatedAt.Unix(), unmarshaledTask.UpdatedAt.Unix())
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		name          string
		taskID        int
		mockSetup     func(sqlmock.Sqlmock)
		expectedTask  *Task
		expectedError error
	}{
		{
			name:   "Successfully get task",
			taskID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				createdAt := time.Now()
				updatedAt := time.Now()
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
					}).AddRow(1, "Test Task", "Test Description", "pending", 1, 0, createdAt, updatedAt))
			},
			expectedTask: &Task{
				ID:          1,
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				UserID:      1,
				Position:    0,
			},
			expectedError: nil,
		},
		{
			name:   "Task not found",
			taskID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedTask:  nil,
			expectedError: sql.ErrNoRows,
		},
		{
			name:   "Database error",
			taskID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedTask:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:   "Scan error",
			taskID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM tasks WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "title", "description", "status", "user_id", "position", "created_at", "updated_at",
					}).AddRow("invalid", "Test Task", "Test Description", "pending", 1, 0, time.Now(), time.Now()))
			},
			expectedTask:  nil,
			expectedError: fmt.Errorf("scan error"), // We just need any error here
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Configure mock expectations
			tt.mockSetup(mock)

			// Execute function
			task, err := GetTask(db, tt.taskID)

			// Assert error
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == sql.ErrNoRows || tt.expectedError == sql.ErrConnDone {
					assert.Equal(t, tt.expectedError, err)
				} else {
					// For scan error, just check that there is an error
					assert.NotNil(t, err)
				}
				// For error cases, task should be zero value
				assert.Equal(t, Task{}, task)
				return
			}

			// Assert no error for successful case
			assert.NoError(t, err)

			// Assert task fields
			assert.Equal(t, tt.expectedTask.ID, task.ID)
			assert.Equal(t, tt.expectedTask.Title, task.Title)
			assert.Equal(t, tt.expectedTask.Description, task.Description)
			assert.Equal(t, tt.expectedTask.Status, task.Status)
			assert.Equal(t, tt.expectedTask.UserID, task.UserID)
			assert.Equal(t, tt.expectedTask.Position, task.Position)
			assert.NotZero(t, task.CreatedAt)
			assert.NotZero(t, task.UpdatedAt)

			// Verify that all expected mock calls were made
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateTask(t *testing.T) {
	// Create a new mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mock.ExpectExec("UPDATE tasks").
		WithArgs("Updated Task", "Updated Description", "completed", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create test task
	task := &Task{
		ID:          1,
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "completed",
	}

	// Call the UpdateTask function
	err = task.UpdateTask(db)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTask(t *testing.T) {
	// Create a new mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mock.ExpectExec("DELETE FROM tasks").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the DeleteTask function
	err = DeleteTask(db, 1)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTaskPosition(t *testing.T) {
	tests := []struct {
		name        string
		task        Task
		userID      int
		newPosition int
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Move task forward",
			task: Task{
				ID:       1,
				Position: 1,
			},
			userID:      1,
			newPosition: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT position FROM tasks WHERE id = \\$1 AND user_id = \\$2").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(1))
				mock.ExpectExec("UPDATE tasks SET position = position - 1, updated_at = \\$1 "+
					"WHERE user_id = \\$2 AND position > \\$3 AND position <= \\$4").
					WithArgs(sqlmock.AnyArg(), 1, 1, 3).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec("UPDATE tasks SET position = \\$1, updated_at = \\$2 "+
					"WHERE id = \\$3 AND user_id = \\$4").
					WithArgs(3, sqlmock.AnyArg(), 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Move task backward",
			task: Task{
				ID:       1,
				Position: 3,
			},
			userID:      1,
			newPosition: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT position FROM tasks WHERE id = \\$1 AND user_id = \\$2").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(3))
				mock.ExpectExec("UPDATE tasks SET position = position \\+ 1, updated_at = \\$1 "+
					"WHERE user_id = \\$2 AND position >= \\$3 AND position < \\$4").
					WithArgs(sqlmock.AnyArg(), 1, 1, 3).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec("UPDATE tasks SET position = \\$1, updated_at = \\$2 "+
					"WHERE id = \\$3 AND user_id = \\$4").
					WithArgs(1, sqlmock.AnyArg(), 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Transaction begin error",
			task: Task{ID: 1},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
		{
			name:   "Position query error",
			task:   Task{ID: 1},
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT position FROM tasks").
					WithArgs(1, 1).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Update other tasks error",
			task: Task{
				ID:       1,
				Position: 1,
			},
			userID:      1,
			newPosition: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT position FROM tasks").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(1))
				mock.ExpectExec("UPDATE tasks SET position").
					WithArgs(sqlmock.AnyArg(), 1, 1, 3).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Update current task error",
			task: Task{
				ID:       1,
				Position: 1,
			},
			userID:      1,
			newPosition: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT position FROM tasks").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(1))
				mock.ExpectExec("UPDATE tasks SET position = position - 1").
					WithArgs(sqlmock.AnyArg(), 1, 1, 3).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec("UPDATE tasks SET position = \\$1").
					WithArgs(3, sqlmock.AnyArg(), 1, 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Commit error",
			task: Task{
				ID:       1,
				Position: 1,
			},
			userID:      1,
			newPosition: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT position FROM tasks").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(1))
				mock.ExpectExec("UPDATE tasks SET position = position - 1").
					WithArgs(sqlmock.AnyArg(), 1, 1, 3).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec("UPDATE tasks SET position = \\$1").
					WithArgs(3, sqlmock.AnyArg(), 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(sql.ErrConnDone)
				// Removed ExpectRollback() as it's not called after commit error
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Configure mock expectations
			tt.mockSetup(mock)

			// Execute function
			err = tt.task.UpdateTaskPosition(db, tt.userID, tt.newPosition)

			// Assert error
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newPosition, tt.task.Position)
				assert.NotZero(t, tt.task.UpdatedAt)
			}

			// Verify that all expected mock calls were made
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name          string
		status        string
		expectedError string
	}{
		{
			name:          "Valid status - pending",
			status:        StatusPending,
			expectedError: "",
		},
		{
			name:          "Valid status - in_progress",
			status:        StatusInProgress,
			expectedError: "",
		},
		{
			name:          "Valid status - completed",
			status:        StatusCompleted,
			expectedError: "",
		},
		{
			name:          "Empty status",
			status:        "",
			expectedError: "invalid status: ",
		},
		{
			name:          "Invalid status",
			status:        "invalid_status",
			expectedError: "invalid status: invalid_status",
		},
		{
			name:          "Case sensitive check",
			status:        "PENDING",
			expectedError: "invalid status: PENDING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Status: tt.status}
			err := task.ValidateStatus()

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			}
		})
	}
}
func TestCreateTask(t *testing.T) {
	tests := []struct {
		name        string
		task        Task
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectedID  int
	}{
		{
			name: "Successful task creation",
			task: Task{
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				UserID:      1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect transaction begin
				mock.ExpectBegin()

				// Expect update of existing tasks positions
				mock.ExpectExec("UPDATE tasks SET position = position \\+ 1").
					WithArgs(1).                              // userID
					WillReturnResult(sqlmock.NewResult(0, 2)) // 2 rows affected

				// Expect task insertion
				mock.ExpectQuery("INSERT INTO tasks \\(title, description, status, user_id, position, created_at, updated_at\\)").
					WithArgs(
						"Test Task",
						"Test Description",
						"pending",
						1,
						0,
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Expect transaction commit
				mock.ExpectCommit()
			},
			expectError: false,
			expectedID:  1,
		},
		{
			name: "Transaction begin error",
			task: Task{
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				UserID:      1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
		{
			name: "Position update error",
			task: Task{
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				UserID:      1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE tasks SET position = position \\+ 1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Insert error",
			task: Task{
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				UserID:      1,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE tasks SET position = position \\+ 1").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("INSERT INTO tasks").
					WithArgs(
						"Test Task",
						"Test Description",
						"pending",
						1,
						0,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Configure mock expectations
			tt.mockSetup(mock)

			// Execute function
			err = tt.task.CreateTask(db)

			// Assert error
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedID, tt.task.ID)
			assert.Equal(t, 0, tt.task.Position)
			assert.NotZero(t, tt.task.CreatedAt)
			assert.NotZero(t, tt.task.UpdatedAt)

			// Verify that all expected mock calls were made
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
