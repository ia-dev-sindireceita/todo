package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// Mock repositories for testing
type mockTaskRepositoryForComplete struct {
	tasks map[string]*application.Task
}

func (m *mockTaskRepositoryForComplete) Create(ctx context.Context, task *application.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForComplete) Update(ctx context.Context, task *application.Task) error {
	if _, exists := m.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForComplete) Delete(ctx context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepositoryForComplete) FindByID(ctx context.Context, id string) (*application.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (m *mockTaskRepositoryForComplete) FindByOwnerID(ctx context.Context, ownerID string) ([]*application.Task, error) {
	var tasks []*application.Task
	for _, task := range m.tasks {
		if task.OwnerID == ownerID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *mockTaskRepositoryForComplete) FindSharedWithUser(ctx context.Context, userID string) ([]*application.Task, error) {
	return []*application.Task{}, nil
}

type mockTaskServiceForComplete struct {
	canAccess bool
	canModify bool
}

func (m *mockTaskServiceForComplete) CanUserAccessTask(ctx context.Context, taskID, userID string) (bool, error) {
	return m.canAccess, nil
}

func (m *mockTaskServiceForComplete) CanUserModifyTask(ctx context.Context, taskID, userID string) (bool, error) {
	return m.canModify, nil
}

func TestCompleteTaskUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		taskID      string
		userID      string
		setupTask   func(*mockTaskRepositoryForComplete)
		canModify   bool
		wantErr     bool
		wantStatus  application.TaskStatus
		errorMsg    string
	}{
		{
			name:   "should complete pending task when user is owner",
			taskID: "task-1",
			userID: "user-1",
			setupTask: func(repo *mockTaskRepositoryForComplete) {
				task, _ := application.NewTask("task-1", "Test Task", "Description", application.StatusPending, "user-1")
				repo.tasks["task-1"] = task
			},
			canModify:  true,
			wantErr:    false,
			wantStatus: application.StatusCompleted,
		},
		{
			name:   "should complete in_progress task when user is owner",
			taskID: "task-2",
			userID: "user-1",
			setupTask: func(repo *mockTaskRepositoryForComplete) {
				task, _ := application.NewTask("task-2", "Test Task", "Description", application.StatusInProgress, "user-1")
				repo.tasks["task-2"] = task
			},
			canModify:  true,
			wantErr:    false,
			wantStatus: application.StatusCompleted,
		},
		{
			name:   "should fail if task not found",
			taskID: "nonexistent",
			userID: "user-1",
			setupTask: func(repo *mockTaskRepositoryForComplete) {
				// No task setup
			},
			canModify: true,
			wantErr:   true,
			errorMsg:  "task not found",
		},
		{
			name:   "should fail if user cannot modify task",
			taskID: "task-3",
			userID: "user-2",
			setupTask: func(repo *mockTaskRepositoryForComplete) {
				task, _ := application.NewTask("task-3", "Test Task", "Description", application.StatusPending, "user-1")
				repo.tasks["task-3"] = task
			},
			canModify: false,
			wantErr:   true,
			errorMsg:  "user does not have permission to modify this task",
		},
		{
			name:   "should fail if task already completed",
			taskID: "task-4",
			userID: "user-1",
			setupTask: func(repo *mockTaskRepositoryForComplete) {
				task, _ := application.NewTask("task-4", "Test Task", "Description", application.StatusCompleted, "user-1")
				repo.tasks["task-4"] = task
			},
			canModify: true,
			wantErr:   true,
			errorMsg:  "task is already completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockTaskRepositoryForComplete{
				tasks: make(map[string]*application.Task),
			}
			tt.setupTask(mockRepo)

			mockService := &mockTaskServiceForComplete{
				canModify: tt.canModify,
			}

			useCase := NewCompleteTaskUseCase(mockRepo, mockService)
			task, err := useCase.Execute(context.Background(), tt.taskID, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error but got nil")
					return
				}
				if task != nil {
					t.Errorf("Execute() expected nil task on error, got %v", task)
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Execute() error = %v, want %v", err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
				return
			}

			if task == nil {
				t.Error("Execute() expected task to be returned, got nil")
				return
			}

			// Verify task status was updated
			if task.Status != tt.wantStatus {
				t.Errorf("Execute() task status = %v, want %v", task.Status, tt.wantStatus)
			}
		})
	}
}
