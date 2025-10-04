package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

// ShareTaskUseCase handles sharing a task with another user
type ShareTaskUseCase struct {
	taskRepo    repository.TaskRepository
	shareRepo   repository.ShareRepository
	taskService *service.TaskService
}

// NewShareTaskUseCase creates a new ShareTaskUseCase
func NewShareTaskUseCase(taskRepo repository.TaskRepository, shareRepo repository.ShareRepository, taskService *service.TaskService) *ShareTaskUseCase {
	return &ShareTaskUseCase{
		taskRepo:    taskRepo,
		shareRepo:   shareRepo,
		taskService: taskService,
	}
}

// Execute shares a task with a user
func (uc *ShareTaskUseCase) Execute(ctx context.Context, taskID, ownerID, shareWithUserID string) error {
	// Check if requesting user is the owner
	canModify, err := uc.taskService.CanUserModifyTask(ctx, taskID, ownerID)
	if err != nil {
		return err
	}
	if !canModify {
		return errors.New("only the task owner can share the task")
	}

	// Cannot share with self
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return err
	}
	if task.OwnerID == shareWithUserID {
		return errors.New("cannot share task with yourself")
	}

	// Share the task
	return uc.shareRepo.Share(ctx, taskID, shareWithUserID)
}
