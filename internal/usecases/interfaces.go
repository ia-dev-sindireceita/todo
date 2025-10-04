package usecases

import (
	"context"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// LoginUseCaseInterface defines the interface for login operations
type LoginUseCaseInterface interface {
	Execute(ctx context.Context, email, password string) (string, error)
}

// RegisterUseCaseInterface defines the interface for registration operations
type RegisterUseCaseInterface interface {
	Execute(ctx context.Context, name, email, password string) (*application.User, error)
}
