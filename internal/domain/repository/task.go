package repository

import (
	"context"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	// Create creates a new task
	Create(ctx context.Context, task *application.Task) error

	// Update updates an existing task
	Update(ctx context.Context, task *application.Task) error

	// Delete deletes a task by ID
	Delete(ctx context.Context, id string) error

	// FindByID finds a task by ID
	FindByID(ctx context.Context, id string) (*application.Task, error)

	// FindByOwnerID finds all tasks owned by a user
	FindByOwnerID(ctx context.Context, ownerID string) ([]*application.Task, error)

	// FindSharedWithUser finds all tasks shared with a user
	FindSharedWithUser(ctx context.Context, userID string) ([]*application.Task, error)
}
