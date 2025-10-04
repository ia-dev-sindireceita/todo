package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
)

// TaskServiceInterface defines the interface for task service operations
type TaskServiceInterface interface {
	CanUserAccessTask(ctx context.Context, taskID, userID string) (bool, error)
	CanUserModifyTask(ctx context.Context, taskID, userID string) (bool, error)
}

// CompleteTaskUseCase handles completing a task
type CompleteTaskUseCase struct {
	taskRepo    repository.TaskRepository
	taskService TaskServiceInterface
}

// NewCompleteTaskUseCase creates a new CompleteTaskUseCase
func NewCompleteTaskUseCase(
	taskRepo repository.TaskRepository,
	taskService TaskServiceInterface,
) *CompleteTaskUseCase {
	return &CompleteTaskUseCase{
		taskRepo:    taskRepo,
		taskService: taskService,
	}
}

// Execute completes a task
func (uc *CompleteTaskUseCase) Execute(ctx context.Context, taskID, userID string) error {
	// Find the task
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return errors.New("task not found")
	}

	// Check if user can modify the task (must be owner)
	canModify, err := uc.taskService.CanUserModifyTask(ctx, taskID, userID)
	if err != nil {
		return err
	}
	if !canModify {
		return errors.New("user does not have permission to modify this task")
	}

	// Complete the task
	if err := task.CompleteTask(); err != nil {
		return err
	}

	// Update in repository
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	return nil
}
