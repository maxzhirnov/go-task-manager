package models

import (
	"log"
	"time"

	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      int       `json:"user_id"` // Associate task with a user
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetTasks(db database.DB, userID int) ([]Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTask(db database.DB, id int) (Task, error) {
	var t Task
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (t *Task) CreateTask(db database.DB) error {
	query := `
        INSERT INTO tasks (title, description, status, user_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	err := db.QueryRow(query, t.Title, t.Description, t.Status, t.UserID, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
	if err != nil {
		log.Printf("Error inserting task into database: %v", err) // Log the error
		return err
	}

	return nil
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

func DeleteTask(db database.DB, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
