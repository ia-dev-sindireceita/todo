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

// CreateTaskUseCaseInterface defines the interface for creating tasks
type CreateTaskUseCaseInterface interface {
	Execute(ctx context.Context, title, description, ownerID string) (*application.Task, error)
}

// GetTaskUseCaseInterface defines the interface for getting a single task
type GetTaskUseCaseInterface interface {
	Execute(ctx context.Context, taskID, userID string) (*application.Task, error)
}

// UpdateTaskUseCaseInterface defines the interface for updating tasks
type UpdateTaskUseCaseInterface interface {
	Execute(ctx context.Context, taskID, title, description string, status application.TaskStatus, userID string) error
}

// DeleteTaskUseCaseInterface defines the interface for deleting tasks
type DeleteTaskUseCaseInterface interface {
	Execute(ctx context.Context, taskID, userID string) error
}

// ListTasksUseCaseInterface defines the interface for listing user's tasks
type ListTasksUseCaseInterface interface {
	Execute(ctx context.Context, userID string) ([]*application.Task, error)
}

// ListSharedTasksUseCaseInterface defines the interface for listing shared tasks
type ListSharedTasksUseCaseInterface interface {
	Execute(ctx context.Context, userID string) ([]*application.Task, error)
}

// CompleteTaskUseCaseInterface defines the interface for completing tasks
type CompleteTaskUseCaseInterface interface {
	Execute(ctx context.Context, taskID, userID string) (*application.Task, error)
}

// ShareTaskUseCaseInterface defines the interface for sharing tasks
type ShareTaskUseCaseInterface interface {
	Execute(ctx context.Context, taskID, ownerID, shareWithUserID string) error
}

// ExportTasksPDFUseCaseInterface defines the interface for exporting tasks to PDF
type ExportTasksPDFUseCaseInterface interface {
	Execute(ctx context.Context, ownerID string) ([]byte, error)
}
