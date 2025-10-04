package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

// DeleteTaskUseCase handles task deletion
type DeleteTaskUseCase struct {
	taskRepo    repository.TaskRepository
	taskService *service.TaskService
}

// NewDeleteTaskUseCase creates a new DeleteTaskUseCase
func NewDeleteTaskUseCase(taskRepo repository.TaskRepository, taskService *service.TaskService) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		taskRepo:    taskRepo,
		taskService: taskService,
	}
}

// Execute deletes a task
func (uc *DeleteTaskUseCase) Execute(ctx context.Context, taskID, userID string) error {
	// Check if user can modify (delete) task
	canModify, err := uc.taskService.CanUserModifyTask(ctx, taskID, userID)
	if err != nil {
		return err
	}
	if !canModify {
		return errors.New("user does not have permission to delete this task")
	}

	// Delete task
	return uc.taskRepo.Delete(ctx, taskID)
}
