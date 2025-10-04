package service

import (
	"context"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
)

// TaskService provides general business logic for tasks
type TaskService struct {
	taskRepo  repository.TaskRepository
	shareRepo repository.ShareRepository
}

// NewTaskService creates a new TaskService
func NewTaskService(taskRepo repository.TaskRepository, shareRepo repository.ShareRepository) *TaskService {
	return &TaskService{
		taskRepo:  taskRepo,
		shareRepo: shareRepo,
	}
}

// CanUserAccessTask checks if a user can access a task (owner or shared with)
func (s *TaskService) CanUserAccessTask(ctx context.Context, taskID, userID string) (bool, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return false, err
	}

	// Owner can always access
	if task.OwnerID == userID {
		return true, nil
	}

	// Check if shared with user
	isShared, err := s.shareRepo.IsSharedWith(ctx, taskID, userID)
	if err != nil {
		return false, err
	}

	return isShared, nil
}

// CanUserModifyTask checks if a user can modify a task (only owner)
func (s *TaskService) CanUserModifyTask(ctx context.Context, taskID, userID string) (bool, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return false, err
	}

	return task.OwnerID == userID, nil
}
