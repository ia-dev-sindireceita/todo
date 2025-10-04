package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// SQLiteUserRepository implements repository.UserRepository using SQLite
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLiteUserRepository
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Create creates a new user using prepared statement
func (r *SQLiteUserRepository) Create(ctx context.Context, user *application.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, created_at)
	          VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
	return err
}

// FindByID finds a user by ID using prepared statement
func (r *SQLiteUserRepository) FindByID(ctx context.Context, id string) (*application.User, error) {
	query := `SELECT id, name, email, password_hash, created_at
	          FROM users WHERE id = ?`

	var user application.User
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &user, nil
}

// FindByEmail finds a user by email using prepared statement
func (r *SQLiteUserRepository) FindByEmail(ctx context.Context, email string) (*application.User, error) {
	query := `SELECT id, name, email, password_hash, created_at
	          FROM users WHERE email = ?`

	var user application.User
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &user, nil
}

// Update updates an existing user using prepared statement
func (r *SQLiteUserRepository) Update(ctx context.Context, user *application.User) error {
	query := `UPDATE users SET name = ?, email = ?, password_hash = ?
	          WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.ID,
	)
	return err
}

// Delete deletes a user using prepared statement
func (r *SQLiteUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
