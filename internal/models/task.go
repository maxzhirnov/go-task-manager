package models

import (
	"fmt"
	"log"
	"time"

	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusDeleted    = "deleted"
)

var ValidStatuses = []string{StatusPending, StatusInProgress, StatusCompleted}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      int       `json:"user_id"` // Associate task with a user
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Position    int       `json:"position"`
}

func GetTasks(db database.DB, userID int) ([]Task, error) {
	query := `SELECT id, title, description, status, user_id, position, created_at, updated_at 
              FROM tasks 
              WHERE user_id = $1 
			  AND status != 'deleted'
              ORDER BY position ASC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status,
			&t.UserID, &t.Position, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTask(db database.DB, id int) (Task, error) {
	var t Task
	query := `SELECT id, title, description, status, user_id, position, created_at, updated_at 
              FROM tasks 
              WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Status,
		&t.UserID, &t.Position, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (t *Task) CreateTask(db database.DB) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Increment positions of existing tasks
	_, err = tx.Exec(`
        UPDATE tasks 
        SET position = position + 1
        WHERE user_id = $1`, t.UserID)
	if err != nil {
		return err
	}

	// Set the position of the new task to 0
	t.Position = 0

	// Insert the new task
	query := `
        INSERT INTO tasks (title, description, status, user_id, position, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	err = tx.QueryRow(query, t.Title, t.Description, t.Status, t.UserID, t.Position, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
	if err != nil {
		log.Printf("Error inserting task into database: %v", err)
		return err
	}

	return tx.Commit()
}

func (t *Task) UpdateTask(db database.DB) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5`

	t.UpdatedAt = time.Now()

	_, err := db.Exec(query, t.Title, t.Description, t.Status, t.UpdatedAt, t.ID)
	return err
}

func (t *Task) UpdateTaskPosition(db database.DB, userID int, newPosition int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get current position
	var oldPosition int
	err = tx.QueryRow(`
        SELECT position 
        FROM tasks 
        WHERE id = $1 AND user_id = $2`, t.ID, userID).Scan(&oldPosition)
	if err != nil {
		return err
	}

	// Update positions of other tasks
	if oldPosition < newPosition {
		_, err = tx.Exec(`
            UPDATE tasks 
            SET position = position - 1,
                updated_at = $1
            WHERE user_id = $2 
            AND position > $3 
            AND position <= $4`,
			time.Now(), userID, oldPosition, newPosition)
	} else {
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
		return err
	}

	// Update position of the current task
	t.Position = newPosition
	t.UpdatedAt = time.Now()
	_, err = tx.Exec(`
        UPDATE tasks 
        SET position = $1,
            updated_at = $2
        WHERE id = $3 AND user_id = $4`,
		t.Position, t.UpdatedAt, t.ID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (t *Task) ValidateStatus() error {
	for _, status := range ValidStatuses {
		if t.Status == status {
			return nil
		}
	}
	return fmt.Errorf("invalid status: %s", t.Status)
}

func DeleteTask(db database.DB, id int) error {
	query := `
        UPDATE tasks 
        SET status = 'deleted', updated_at = $1 
        WHERE id = $2
    `
	result, err := db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
