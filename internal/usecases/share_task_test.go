package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

func TestShareTaskUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	taskID := "task-1"
	ownerID := "user-1"
	shareWithUserID := "user-2"

	task, _ := application.NewTask(taskID, "Test Task", "Description", application.StatusPending, ownerID)

	taskRepo := &mockTaskRepositoryForShare{
		tasks: map[string]*application.Task{
			taskID: task,
		},
	}
	shareRepo := &mockShareRepositoryForShare{}
	taskService := service.NewTaskService(taskRepo, shareRepo)

	useCase := NewShareTaskUseCase(taskRepo, shareRepo, taskService)

	err := useCase.Execute(ctx, taskID, ownerID, shareWithUserID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify task was shared
	if !shareRepo.shared {
		t.Error("Expected task to be shared")
	}
}

func TestShareTaskUseCase_Execute_OnlyOwnerCanShare(t *testing.T) {
	ctx := context.Background()
	taskID := "task-1"
	ownerID := "user-1"
	nonOwnerID := "user-2"
	shareWithUserID := "user-3"

	task, _ := application.NewTask(taskID, "Test Task", "Description", application.StatusPending, ownerID)

	taskRepo := &mockTaskRepositoryForShare{
		tasks: map[string]*application.Task{
			taskID: task,
		},
	}
	shareRepo := &mockShareRepositoryForShare{}
	taskService := service.NewTaskService(taskRepo, shareRepo)

	useCase := NewShareTaskUseCase(taskRepo, shareRepo, taskService)

	// Non-owner tries to share
	err := useCase.Execute(ctx, taskID, nonOwnerID, shareWithUserID)
	if err == nil {
		t.Error("Expected error when non-owner tries to share")
	}
	if err.Error() != "only the task owner can share the task" {
		t.Errorf("Expected 'only the task owner can share the task' error, got %v", err)
	}

	// Verify task was NOT shared
	if shareRepo.shared {
		t.Error("Expected task NOT to be shared")
	}
}

func TestShareTaskUseCase_Execute_CannotShareWithSelf(t *testing.T) {
	ctx := context.Background()
	taskID := "task-1"
	ownerID := "user-1"

	task, _ := application.NewTask(taskID, "Test Task", "Description", application.StatusPending, ownerID)

	taskRepo := &mockTaskRepositoryForShare{
		tasks: map[string]*application.Task{
			taskID: task,
		},
	}
	shareRepo := &mockShareRepositoryForShare{}
	taskService := service.NewTaskService(taskRepo, shareRepo)

	useCase := NewShareTaskUseCase(taskRepo, shareRepo, taskService)

	// Try to share with self
	err := useCase.Execute(ctx, taskID, ownerID, ownerID)
	if err == nil {
		t.Error("Expected error when sharing with self")
	}
	if err.Error() != "cannot share task with yourself" {
		t.Errorf("Expected 'cannot share task with yourself' error, got %v", err)
	}

	// Verify task was NOT shared
	if shareRepo.shared {
		t.Error("Expected task NOT to be shared")
	}
}

func TestShareTaskUseCase_Execute_TaskNotFound(t *testing.T) {
	ctx := context.Background()
	taskID := "nonexistent"
	ownerID := "user-1"
	shareWithUserID := "user-2"

	taskRepo := &mockTaskRepositoryForShare{
		tasks: map[string]*application.Task{},
	}
	shareRepo := &mockShareRepositoryForShare{}
	taskService := service.NewTaskService(taskRepo, shareRepo)

	useCase := NewShareTaskUseCase(taskRepo, shareRepo, taskService)

	err := useCase.Execute(ctx, taskID, ownerID, shareWithUserID)
	if err == nil {
		t.Error("Expected error when task not found")
	}
}

// Mock repositories for testing
type mockTaskRepositoryForShare struct {
	tasks map[string]*application.Task
}

func (m *mockTaskRepositoryForShare) Create(ctx context.Context, task *application.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForShare) Update(ctx context.Context, task *application.Task) error {
	if _, exists := m.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForShare) Delete(ctx context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepositoryForShare) FindByID(ctx context.Context, id string) (*application.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (m *mockTaskRepositoryForShare) FindByOwnerID(ctx context.Context, ownerID string) ([]*application.Task, error) {
	var tasks []*application.Task
	for _, task := range m.tasks {
		if task.OwnerID == ownerID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *mockTaskRepositoryForShare) FindSharedWithUser(ctx context.Context, userID string) ([]*application.Task, error) {
	return []*application.Task{}, nil
}

type mockShareRepositoryForShare struct {
	shared bool
	shares map[string][]string
}

func (m *mockShareRepositoryForShare) Share(ctx context.Context, taskID, userID string) error {
	m.shared = true
	if m.shares == nil {
		m.shares = make(map[string][]string)
	}
	m.shares[taskID] = append(m.shares[taskID], userID)
	return nil
}

func (m *mockShareRepositoryForShare) Unshare(ctx context.Context, taskID, userID string) error {
	return nil
}

func (m *mockShareRepositoryForShare) FindSharedUsers(ctx context.Context, taskID string) ([]string, error) {
	if users, ok := m.shares[taskID]; ok {
		return users, nil
	}
	return []string{}, nil
}

func (m *mockShareRepositoryForShare) IsSharedWith(ctx context.Context, taskID, userID string) (bool, error) {
	if users, ok := m.shares[taskID]; ok {
		for _, u := range users {
			if u == userID {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *mockShareRepositoryForShare) DeleteAllShares(ctx context.Context, taskID string) error {
	delete(m.shares, taskID)
	return nil
}
