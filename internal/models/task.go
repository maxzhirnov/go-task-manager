// Package models provides data structures and database operations
// for the task management application.
package models

import (
	"fmt"
	"log"
	"time"

	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

// Task status constants define all possible states of a task
const (
	// StatusPending represents a task that hasn't been started
	StatusPending = "pending"

	// StatusInProgress represents a task that is currently being worked on
	StatusInProgress = "in_progress"

	// StatusCompleted represents a task that has been finished
	StatusCompleted = "completed"

	// StatusDeleted represents a soft-deleted task
	StatusDeleted = "deleted"
)

// ValidStatuses defines the list of allowed task statuses for validation.
// Note: StatusDeleted is not included as it's only set internally.
var ValidStatuses = []string{StatusPending, StatusInProgress, StatusCompleted}

// Task represents a single task in the system.
// It contains all task-related information including its current state,
// position in the user's task list, and timestamps.
type Task struct {
	// ID uniquely identifies the task
	ID int `json:"id"`

	// Title is the main identifier of the task
	Title string `json:"title"`

	// Description provides detailed information about the task
	Description string `json:"description"`

	// Status represents the current state of the task
	// Must be one of ValidStatuses
	Status string `json:"status"`

	// UserID associates the task with a specific user
	UserID int `json:"user_id"`

	// CreatedAt stores the timestamp when the task was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt stores the timestamp of the last modification
	UpdatedAt time.Time `json:"updated_at"`

	// Position represents the task's order in the user's task list
	Position int `json:"position"`
}

// GetTasks retrieves all active tasks for a specific user.
//
// It returns tasks ordered by their position, excluding soft-deleted tasks.
// The function performs a database query to fetch tasks associated with the
// provided user ID.
//
// Parameters:
//   - db: Database interface for executing queries
//   - userID: The ID of the user whose tasks to retrieve
//
// Returns:
//   - []Task: Slice of tasks belonging to the user
//   - error: Database error if query fails
//
// Query Details:
//   - Excludes tasks with status 'deleted'
//   - Orders tasks by position ascending
//   - Includes all task fields
//
// Example Usage:
//
//	tasks, err := GetTasks(db, userID)
//	if err != nil {
//	    return fmt.Errorf("failed to fetch tasks: %w", err)
//	}
func GetTasks(db database.DB, userID int) ([]Task, error) {
	// SQL query to fetch active tasks for user
	query := `SELECT id, title, description, status, user_id, position, created_at, updated_at 
              FROM tasks 
              WHERE user_id = $1 
              AND status != 'deleted'
              ORDER BY position ASC`

	// Execute query with user ID
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through results and build tasks slice
	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.UserID,
			&t.Position,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// GetTask retrieves a single task by its ID.
//
// It performs a database query to fetch a specific task's details.
// Unlike GetTasks, this function returns deleted tasks as well,
// allowing for administrative access and restoration capabilities.
//
// Parameters:
//   - db: Database interface for executing queries
//   - id: The unique identifier of the task to retrieve
//
// Returns:
//   - Task: The requested task's data
//   - error: Database error or sql.ErrNoRows if task not found
//
// Query Details:
//   - Retrieves all task fields
//   - Includes deleted tasks
//   - Uses direct ID matching
//
// Example Usage:
//
//	task, err := GetTask(db, taskID)
//	if err == sql.ErrNoRows {
//	    return Task{}, fmt.Errorf("task not found")
//	}
//	if err != nil {
//	    return Task{}, fmt.Errorf("failed to fetch task: %w", err)
//	}
func GetTask(db database.DB, id int) (Task, error) {
	var t Task

	// SQL query to fetch task by ID
	query := `SELECT id, title, description, status, user_id, position, created_at, updated_at 
              FROM tasks 
              WHERE id = $1`

	// Execute query and scan results into Task struct
	err := db.QueryRow(query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.UserID,
		&t.Position,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	return t, err
}

// CreateTask inserts a new task into the database and updates task positions.
//
// This method uses a transaction to ensure atomicity of the operation:
// 1. Increments positions of existing tasks for the user
// 2. Inserts the new task at position 0
//
// The method also sets the creation and update timestamps.
//
// Parameters:
//   - db: Database interface for executing queries
//
// Returns:
//   - error: Any error encountered during the process
//
// Side Effects:
//   - Sets t.ID with the newly created task's ID
//   - Sets t.CreatedAt and t.UpdatedAt to current time
//   - Sets t.Position to 0
//
// Example Usage:
//
//	task := &Task{
//	    Title:       "New Task",
//	    Description: "Task details",
//	    Status:      StatusPending,
//	    UserID:      userID,
//	}
//	if err := task.CreateTask(db); err != nil {
//	    return fmt.Errorf("failed to create task: %w", err)
//	}
func (t *Task) CreateTask(db database.DB) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback in case of error

	// Increment positions of existing tasks
	_, err = tx.Exec(`
        UPDATE tasks 
        SET position = position + 1
        WHERE user_id = $1`, t.UserID)
	if err != nil {
		return fmt.Errorf("failed to update task positions: %w", err)
	}

	// Set the position of the new task to 0
	t.Position = 0

	// Set creation and update timestamps
	t.CreatedAt = time.Now()
	t.UpdatedAt = t.CreatedAt

	// Insert the new task
	query := `
        INSERT INTO tasks (title, description, status, user_id, position, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	err = tx.QueryRow(query, t.Title, t.Description, t.Status, t.UserID, t.Position, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
	if err != nil {
		log.Printf("Error inserting task into database: %v", err)
		return fmt.Errorf("failed to insert task: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateTask modifies an existing task's details in the database.
//
// This method updates the task's title, description, and status,
// while automatically updating the updated_at timestamp. The task's
// ID must be set before calling this method.
//
// Parameters:
//   - db: Database interface for executing queries
//
// Returns:
//   - error: Database error if update fails
//
// Fields Updated:
//   - title
//   - description
//   - status
//   - updated_at (automatically set to current time)
//
// Example Usage:
//
//	task.Title = "Updated Title"
//	task.Status = StatusInProgress
//	if err := task.UpdateTask(db); err != nil {
//	    return fmt.Errorf("failed to update task: %w", err)
//	}
//
// Note: This method does not update the task's position or user_id
// as these should be modified through separate specialized methods.
func (t *Task) UpdateTask(db database.DB) error {
	// SQL query to update task fields
	query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, updated_at = $4
        WHERE id = $5`

	// Set current timestamp
	t.UpdatedAt = time.Now()

	// Execute update query
	_, err := db.Exec(query, t.Title, t.Description, t.Status, t.UpdatedAt, t.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// UpdateTaskPosition changes a task's position in the user's task list.
//
// This method uses a transaction to ensure atomicity when reordering tasks.
// It handles both moving a task up (to a lower position number) and down
// (to a higher position number) while maintaining the integrity of all
// task positions.
//
// The process:
// 1. Verifies task ownership and gets current position
// 2. Shifts other tasks' positions to make space
// 3. Updates the target task's position
//
// Parameters:
//   - db: Database interface for executing queries
//   - userID: ID of the task owner (for verification)
//   - newPosition: Desired position for the task
//
// Returns:
//   - error: Any error encountered during the process
//
// Example Usage:
//
//	// Move task to position 3
//	if err := task.UpdateTaskPosition(db, userID, 3); err != nil {
//	    return fmt.Errorf("failed to update task position: %w", err)
//	}
//
// Note: This method updates the updated_at timestamp for all affected tasks.
func (t *Task) UpdateTaskPosition(db database.DB, userID int, newPosition int) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback in case of error

	// Get current position and verify ownership
	var oldPosition int
	err = tx.QueryRow(`
        SELECT position 
        FROM tasks 
        WHERE id = $1 AND user_id = $2`, t.ID, userID).Scan(&oldPosition)
	if err != nil {
		return fmt.Errorf("failed to get current position: %w", err)
	}

	// Update positions of affected tasks
	if oldPosition < newPosition {
		// Moving task down: shift intermediate tasks up
		_, err = tx.Exec(`
            UPDATE tasks 
            SET position = position - 1,
                updated_at = $1
            WHERE user_id = $2 
            AND position > $3 
            AND position <= $4`,
			time.Now(), userID, oldPosition, newPosition)
	} else {
		// Moving task up: shift intermediate tasks down
		_, err = tx.Exec(`
            UPDATE tasks 
            SET position = position + 1,
                updated_at = $1
            WHERE user_id = $2 
            AND position >= $3 
            AND position < $4`,
			time.Now(), userID, newPosition, oldPosition)
	}
	if err != nil {
		return fmt.Errorf("failed to update intermediate positions: %w", err)
	}

	// Update target task's position
	t.Position = newPosition
	t.UpdatedAt = time.Now()
	_, err = tx.Exec(`
        UPDATE tasks 
        SET position = $1,
            updated_at = $2
        WHERE id = $3 AND user_id = $4`,
		t.Position, t.UpdatedAt, t.ID, userID)
	if err != nil {
		return fmt.Errorf("failed to update task position: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ValidateStatus checks if the task's status is one of the allowed values.
//
// It compares the task's status against the ValidStatuses slice to ensure
// only permitted status values are used. This method is typically called
// before creating or updating a task.
//
// Returns:
//   - nil: If the status is valid
//   - error: If the status is not in ValidStatuses
//
// Valid Statuses:
//   - "pending"
//   - "in_progress"
//   - "completed"
//
// Example Usage:
//
//	task := &Task{Status: "invalid_status"}
//	if err := task.ValidateStatus(); err != nil {
//	    return fmt.Errorf("validation failed: %w", err)
//	}
//
// Note: StatusDeleted is not included in valid statuses as it should
// only be set through internal deletion processes.
func (t *Task) ValidateStatus() error {
	// Check if status matches any valid status
	for _, status := range ValidStatuses {
		if t.Status == status {
			return nil
		}
	}

	// Return error for invalid status
	return fmt.Errorf("invalid status: %s", t.Status)
}

// DeleteTask performs a soft delete of a task by marking its status as 'deleted'.
//
// Instead of removing the task from the database, this function updates the
// task's status to 'deleted' and its updated_at timestamp. This approach
// allows for potential task restoration and maintains data history.
//
// Parameters:
//   - db: Database interface for executing queries
//   - id: The unique identifier of the task to delete
//
// Returns:
//   - error: Database error or "task not found" if task doesn't exist
//
// Example Usage:
//
//	if err := DeleteTask(db, taskID); err != nil {
//	    if err.Error() == "task not found" {
//	        return fmt.Errorf("cannot delete: task %d not found", taskID)
//	    }
//	    return fmt.Errorf("failed to delete task: %w", err)
//	}
//
// Note: This is a soft delete operation. The task will still exist in the
// database but won't appear in normal task listings.
func DeleteTask(db database.DB, id int) error {
	// SQL query for soft delete
	query := `
        UPDATE tasks 
        SET status = 'deleted', updated_at = $1 
        WHERE id = $2
    `

	// Execute update
	result, err := db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	// Check if task existed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Return error if no task was found
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
