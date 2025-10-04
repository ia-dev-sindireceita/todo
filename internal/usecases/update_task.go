package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

// UpdateTaskUseCase handles task updates
type UpdateTaskUseCase struct {
	taskRepo    repository.TaskRepository
	taskService *service.TaskService
}

// NewUpdateTaskUseCase creates a new UpdateTaskUseCase
func NewUpdateTaskUseCase(taskRepo repository.TaskRepository, taskService *service.TaskService) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		taskRepo:    taskRepo,
		taskService: taskService,
	}
}

// Execute updates a task
func (uc *UpdateTaskUseCase) Execute(ctx context.Context, taskID, title, description string, status application.TaskStatus, imagePath, userID string) error {
	// Check if user can modify task
	canModify, err := uc.taskService.CanUserModifyTask(ctx, taskID, userID)
	if err != nil {
		return err
	}
	if !canModify {
		return errors.New("user does not have permission to modify this task")
	}

	// Get task
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	// Update task with validation
	if err := task.Update(title, description, status, imagePath); err != nil {
		return err
	}

	// Persist changes
	return uc.taskRepo.Update(ctx, task)
}
