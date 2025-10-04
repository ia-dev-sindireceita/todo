package usecases

import (
	"context"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
)

// ListTasksUseCase handles listing tasks owned by a user
type ListTasksUseCase struct {
	taskRepo repository.TaskRepository
}

// NewListTasksUseCase creates a new ListTasksUseCase
func NewListTasksUseCase(taskRepo repository.TaskRepository) *ListTasksUseCase {
	return &ListTasksUseCase{
		taskRepo: taskRepo,
	}
}

// Execute lists all tasks owned by a user
func (uc *ListTasksUseCase) Execute(ctx context.Context, ownerID string) ([]*application.Task, error) {
	return uc.taskRepo.FindByOwnerID(ctx, ownerID)
}
