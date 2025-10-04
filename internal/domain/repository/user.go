package repository

import (
	"context"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *application.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*application.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*application.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *application.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error
}
