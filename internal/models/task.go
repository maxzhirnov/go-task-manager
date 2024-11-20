package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetTasks(db *sql.DB) ([]Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at DESC`
	rows, err := db.Query(query)
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

func GetTask(db *sql.DB, id int) (Task, error) {
	var t Task
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (t *Task) CreateTask(db *sql.DB) error {
	query := `
        INSERT INTO tasks (title, description, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	return db.QueryRow(query, t.Title, t.Description, t.Status, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
}

func (t *Task) UpdateTask(db *sql.DB) error {
	query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, updated_at = $4
        WHERE id = $5`

	t.UpdatedAt = time.Now()

	_, err := db.Exec(query, t.Title, t.Description, t.Status, t.UpdatedAt, t.ID)
	return err
}

func DeleteTask(db *sql.DB, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
