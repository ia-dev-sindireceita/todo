package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
)

// DeleteTaskImageUseCase handles deleting an image from a task
type DeleteTaskImageUseCase struct {
	taskRepo    repository.TaskRepository
	taskService TaskServiceInterface
}

// NewDeleteTaskImageUseCase creates a new DeleteTaskImageUseCase
func NewDeleteTaskImageUseCase(
	taskRepo repository.TaskRepository,
	taskService TaskServiceInterface,
) *DeleteTaskImageUseCase {
	return &DeleteTaskImageUseCase{
		taskRepo:    taskRepo,
		taskService: taskService,
	}
}

// Execute deletes an image from a task and returns the old image path for cleanup
func (uc *DeleteTaskImageUseCase) Execute(ctx context.Context, taskID, userID string) (string, error) {
	// Find the task
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return "", errors.New("task not found")
	}

	// Check if user can modify the task (must be owner)
	canModify, err := uc.taskService.CanUserModifyTask(ctx, taskID, userID)
	if err != nil {
		return "", err
	}
	if !canModify {
		return "", errors.New("user does not have permission to modify this task")
	}

	// Store old image path for cleanup
	oldImagePath := task.ImagePath

	// Remove the image from the task
	if err := task.RemoveImage(); err != nil {
		return "", err
	}

	// Update in repository
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return "", err
	}

	return oldImagePath, nil
}
