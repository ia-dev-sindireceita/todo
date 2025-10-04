package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
)

// ReplaceTaskImageUseCase handles replacing an image in a task
type ReplaceTaskImageUseCase struct {
	taskRepo    repository.TaskRepository
	taskService TaskServiceInterface
}

// NewReplaceTaskImageUseCase creates a new ReplaceTaskImageUseCase
func NewReplaceTaskImageUseCase(
	taskRepo repository.TaskRepository,
	taskService TaskServiceInterface,
) *ReplaceTaskImageUseCase {
	return &ReplaceTaskImageUseCase{
		taskRepo:    taskRepo,
		taskService: taskService,
	}
}

// Execute replaces an image in a task and returns the old image path for cleanup
func (uc *ReplaceTaskImageUseCase) Execute(ctx context.Context, taskID, userID, newImagePath string) (string, error) {
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

	// Replace the image in the task
	if err := task.ReplaceImage(newImagePath); err != nil {
		return "", err
	}

	// Update in repository
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		return "", err
	}

	return oldImagePath, nil
}
