package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// Mock TaskRepository
type mockTaskRepository struct {
	tasks         map[string]*application.Task
	findByIDError error
}

func (m *mockTaskRepository) Create(ctx context.Context, task *application.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) Update(ctx context.Context, task *application.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) Delete(ctx context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepository) FindByID(ctx context.Context, id string) (*application.Task, error) {
	if m.findByIDError != nil {
		return nil, m.findByIDError
	}
	task, ok := m.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (m *mockTaskRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*application.Task, error) {
	var tasks []*application.Task
	for _, task := range m.tasks {
		if task.OwnerID == ownerID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *mockTaskRepository) FindSharedWithUser(ctx context.Context, userID string) ([]*application.Task, error) {
	return nil, nil
}

func TestTaskService_CanUserAccessTask(t *testing.T) {
	task, _ := application.NewTask("task-1", "Test Task", "Description", application.StatusPending, "user-1", "")

	mockRepo := &mockTaskRepository{
		tasks: map[string]*application.Task{
			"task-1": task,
		},
	}

	mockShareRepo := &mockShareRepository{
		shares: map[string][]string{},
	}

	service := NewTaskService(mockRepo, mockShareRepo)

	tests := []struct {
		name    string
		taskID  string
		userID  string
		want    bool
		wantErr bool
	}{
		{
			name:    "owner can access",
			taskID:  "task-1",
			userID:  "user-1",
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-owner cannot access",
			taskID:  "task-1",
			userID:  "user-2",
			want:    false,
			wantErr: false,
		},
		{
			name:    "task not found",
			taskID:  "task-999",
			userID:  "user-1",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CanUserAccessTask(context.Background(), tt.taskID, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CanUserAccessTask() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CanUserAccessTask() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("CanUserAccessTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Mock ShareRepository
type mockShareRepository struct {
	shares map[string][]string
}

func (m *mockShareRepository) Share(ctx context.Context, taskID, userID string) error {
	m.shares[taskID] = append(m.shares[taskID], userID)
	return nil
}

func (m *mockShareRepository) Unshare(ctx context.Context, taskID, userID string) error {
	users := m.shares[taskID]
	for i, u := range users {
		if u == userID {
			m.shares[taskID] = append(users[:i], users[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockShareRepository) FindSharedUsers(ctx context.Context, taskID string) ([]string, error) {
	return m.shares[taskID], nil
}

func (m *mockShareRepository) IsSharedWith(ctx context.Context, taskID, userID string) (bool, error) {
	users, ok := m.shares[taskID]
	if !ok {
		return false, nil
	}
	for _, u := range users {
		if u == userID {
			return true, nil
		}
	}
	return false, nil
}

func TestTaskService_CanUserModifyTask(t *testing.T) {
	task, _ := application.NewTask("task-1", "Test Task", "Description", application.StatusPending, "user-1", "")

	mockRepo := &mockTaskRepository{
		tasks: map[string]*application.Task{
			"task-1": task,
		},
	}

	mockShareRepo := &mockShareRepository{
		shares: map[string][]string{},
	}

	service := NewTaskService(mockRepo, mockShareRepo)

	tests := []struct {
		name    string
		taskID  string
		userID  string
		want    bool
		wantErr bool
	}{
		{
			name:    "owner can modify",
			taskID:  "task-1",
			userID:  "user-1",
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-owner cannot modify",
			taskID:  "task-1",
			userID:  "user-2",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CanUserModifyTask(context.Background(), tt.taskID, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CanUserModifyTask() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CanUserModifyTask() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("CanUserModifyTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
