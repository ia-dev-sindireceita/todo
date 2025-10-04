package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

// GetTaskUseCase handles retrieving a single task
type GetTaskUseCase struct {
	taskRepo    repository.TaskRepository
	taskService *service.TaskService
}

// NewGetTaskUseCase creates a new GetTaskUseCase
func NewGetTaskUseCase(taskRepo repository.TaskRepository, taskService *service.TaskService) *GetTaskUseCase {
	return &GetTaskUseCase{
		taskRepo:    taskRepo,
		taskService: taskService,
	}
}

// Execute retrieves a task
func (uc *GetTaskUseCase) Execute(ctx context.Context, taskID, userID string) (*application.Task, error) {
	// Check if user can access task
	canAccess, err := uc.taskService.CanUserAccessTask(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}
	if !canAccess {
		return nil, errors.New("user does not have permission to access this task")
	}

	return uc.taskRepo.FindByID(ctx, taskID)
}
