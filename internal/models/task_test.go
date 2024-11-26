package models

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// func TestGetTasks(t *testing.T) {
// 	// Create a new mock DB
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Mock the expected SQL query and result
// 	mockRows := sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
// 		AddRow(1, "Test Task", "Test Description", "pending", time.Now(), time.Now())
// 	mock.ExpectQuery("SELECT id, title, description, status, created_at, updated_at FROM tasks").
// 		WillReturnRows(mockRows)

// 	// Call the GetTasks function
// 	tasks, err := GetTasks(db, 1)
// 	assert.NoError(t, err)
// 	assert.Len(t, tasks, 1)

// 	// Validate values
// 	assert.Equal(t, 1, tasks[0].ID)
// 	assert.Equal(t, "Test Task", tasks[0].Title)
// 	assert.Equal(t, "Test Description", tasks[0].Description)
// 	assert.Equal(t, "pending", tasks[0].Status)

// 	// Ensure all expectations were met
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestGetTask(t *testing.T) {
// 	// Create a new mock DB
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Mock the expected SQL query and result
// 	mockRow := sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
// 		AddRow(1, "Test Task", "Test Description", "pending", time.Now(), time.Now())
// 	mock.ExpectQuery("SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?").
// 		WithArgs(1).
// 		WillReturnRows(mockRow)

// 	// Call the GetTask function
// 	task, err := GetTask(db, 1)
// 	assert.NoError(t, err)

// 	// Validate values
// 	assert.Equal(t, 1, task.ID)
// 	assert.Equal(t, "Test Task", task.Title)
// 	assert.Equal(t, "Test Description", task.Description)
// 	assert.Equal(t, "pending", task.Status)

// 	// Ensure all expectations were met
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

func TestCreateTask(t *testing.T) {
	// Create a new mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query and result
	mock.ExpectQuery("INSERT INTO tasks").
		WithArgs("Test Task", "Test Description", "pending", 1, sqlmock.AnyArg(), sqlmock.AnyArg()). // Include user_id = 1
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create test task
	task := &Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		UserID:      1, // Set the user ID
	}

	// Call the CreateTask function
	err = task.CreateTask(db)
	assert.NoError(t, err)

	// Validate the task ID
	assert.Equal(t, 1, task.ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
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
