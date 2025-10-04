package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// SQLiteTaskRepository implements repository.TaskRepository using SQLite
type SQLiteTaskRepository struct {
	db *sql.DB
}

// NewSQLiteTaskRepository creates a new SQLiteTaskRepository
func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{db: db}
}

// Create creates a new task using prepared statement
func (r *SQLiteTaskRepository) Create(ctx context.Context, task *application.Task) error {
	query := `INSERT INTO tasks (id, title, description, status, owner_id, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		string(task.Status),
		task.OwnerID,
		task.CreatedAt,
		task.UpdatedAt,
	)
	return err
}

// Update updates an existing task using prepared statement
func (r *SQLiteTaskRepository) Update(ctx context.Context, task *application.Task) error {
	query := `UPDATE tasks SET title = ?, description = ?, status = ?, updated_at = ?
	          WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		task.Title,
		task.Description,
		string(task.Status),
		task.UpdatedAt,
		task.ID,
	)
	return err
}

// Delete deletes a task using prepared statement
func (r *SQLiteTaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// FindByID finds a task by ID using prepared statement
func (r *SQLiteTaskRepository) FindByID(ctx context.Context, id string) (*application.Task, error) {
	query := `SELECT id, title, description, status, owner_id, created_at, updated_at
	          FROM tasks WHERE id = ?`

	var task application.Task
	var status string
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.OwnerID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	task.Status = application.TaskStatus(status)
	task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &task, nil
}

// FindByOwnerID finds all tasks owned by a user using prepared statement
func (r *SQLiteTaskRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*application.Task, error) {
	query := `SELECT id, title, description, status, owner_id, created_at, updated_at
	          FROM tasks WHERE owner_id = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*application.Task
	for rows.Next() {
		var task application.Task
		var status string
		var createdAt, updatedAt string

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&status,
			&task.OwnerID,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		task.Status = application.TaskStatus(status)
		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

// FindSharedWithUser finds all tasks shared with a user using prepared statement
func (r *SQLiteTaskRepository) FindSharedWithUser(ctx context.Context, userID string) ([]*application.Task, error) {
	query := `SELECT t.id, t.title, t.description, t.status, t.owner_id, t.created_at, t.updated_at
	          FROM tasks t
	          INNER JOIN task_shares ts ON t.id = ts.task_id
	          WHERE ts.user_id = ?
	          ORDER BY t.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*application.Task
	for rows.Next() {
		var task application.Task
		var status string
		var createdAt, updatedAt string

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&status,
			&task.OwnerID,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		task.Status = application.TaskStatus(status)
		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}
