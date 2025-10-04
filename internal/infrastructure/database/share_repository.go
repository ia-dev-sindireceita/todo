package database

import (
	"context"
	"database/sql"
)

// SQLiteShareRepository implements repository.ShareRepository using SQLite
type SQLiteShareRepository struct {
	db *sql.DB
}

// NewSQLiteShareRepository creates a new SQLiteShareRepository
func NewSQLiteShareRepository(db *sql.DB) *SQLiteShareRepository {
	return &SQLiteShareRepository{db: db}
}

// Share shares a task with a user using prepared statement
func (r *SQLiteShareRepository) Share(ctx context.Context, taskID, userID string) error {
	query := `INSERT INTO task_shares (task_id, user_id) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, taskID, userID)
	return err
}

// Unshare removes sharing of a task with a user using prepared statement
func (r *SQLiteShareRepository) Unshare(ctx context.Context, taskID, userID string) error {
	query := `DELETE FROM task_shares WHERE task_id = ? AND user_id = ?`
	_, err := r.db.ExecContext(ctx, query, taskID, userID)
	return err
}

// FindSharedUsers finds all users a task is shared with using prepared statement
func (r *SQLiteShareRepository) FindSharedUsers(ctx context.Context, taskID string) ([]string, error) {
	query := `SELECT user_id FROM task_shares WHERE task_id = ?`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, rows.Err()
}

// IsSharedWith checks if a task is shared with a user using prepared statement
func (r *SQLiteShareRepository) IsSharedWith(ctx context.Context, taskID, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM task_shares WHERE task_id = ? AND user_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, taskID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
