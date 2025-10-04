package repository

import (
	"context"
)

// TaskShare represents a task sharing relationship
type TaskShare struct {
	TaskID string
	UserID string
}

// ShareRepository defines the interface for task sharing persistence
type ShareRepository interface {
	// Share shares a task with a user
	Share(ctx context.Context, taskID, userID string) error

	// Unshare removes sharing of a task with a user
	Unshare(ctx context.Context, taskID, userID string) error

	// FindSharedUsers finds all users a task is shared with
	FindSharedUsers(ctx context.Context, taskID string) ([]string, error)

	// IsSharedWith checks if a task is shared with a user
	IsSharedWith(ctx context.Context, taskID, userID string) (bool, error)
}
