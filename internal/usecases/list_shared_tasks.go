package usecases

import (
	"context"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
)

// ListSharedTasksUseCase handles listing tasks shared with a user
type ListSharedTasksUseCase struct {
	taskRepo repository.TaskRepository
}

// NewListSharedTasksUseCase creates a new ListSharedTasksUseCase
func NewListSharedTasksUseCase(taskRepo repository.TaskRepository) *ListSharedTasksUseCase {
	return &ListSharedTasksUseCase{
		taskRepo: taskRepo,
	}
}

// Execute lists all tasks shared with a user
func (uc *ListSharedTasksUseCase) Execute(ctx context.Context, userID string) ([]*application.Task, error) {
	return uc.taskRepo.FindSharedWithUser(ctx, userID)
}
