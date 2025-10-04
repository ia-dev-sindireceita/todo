package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
)

// CreateTaskUseCase handles task creation
type CreateTaskUseCase struct {
	taskRepo repository.TaskRepository
}

// NewCreateTaskUseCase creates a new CreateTaskUseCase
func NewCreateTaskUseCase(taskRepo repository.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		taskRepo: taskRepo,
	}
}

// Execute creates a new task
func (uc *CreateTaskUseCase) Execute(ctx context.Context, title, description, ownerID, imagePath string) (*application.Task, error) {
	// Generate unique ID
	id := uuid.New().String()

	// Create task entity with validation
	task, err := application.NewTask(id, title, description, application.StatusPending, ownerID, imagePath)
	if err != nil {
		return nil, err
	}

	// Persist task
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}
