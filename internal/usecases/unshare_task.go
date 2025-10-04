package usecases

import (
	"context"
	"errors"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

// UnshareTaskUseCase handles removing task sharing
type UnshareTaskUseCase struct {
	shareRepo   repository.ShareRepository
	taskService *service.TaskService
}

// NewUnshareTaskUseCase creates a new UnshareTaskUseCase
func NewUnshareTaskUseCase(shareRepo repository.ShareRepository, taskService *service.TaskService) *UnshareTaskUseCase {
	return &UnshareTaskUseCase{
		shareRepo:   shareRepo,
		taskService: taskService,
	}
}

// Execute removes sharing of a task
func (uc *UnshareTaskUseCase) Execute(ctx context.Context, taskID, ownerID, userID string) error {
	// Check if requesting user is the owner
	canModify, err := uc.taskService.CanUserModifyTask(ctx, taskID, ownerID)
	if err != nil {
		return err
	}
	if !canModify {
		return errors.New("only the task owner can unshare the task")
	}

	return uc.shareRepo.Unshare(ctx, taskID, userID)
}
