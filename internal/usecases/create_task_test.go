package usecases

import (
	"context"
	"testing"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

func TestCreateTaskUseCase_Execute(t *testing.T) {
	mockRepo := &mockTaskRepository{
		tasks: make(map[string]*application.Task),
	}

	useCase := NewCreateTaskUseCase(mockRepo)

	tests := []struct {
		name        string
		title       string
		description string
		ownerID     string
		imagePath   string
		wantErr     bool
	}{
		{
			name:        "valid task creation",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     false,
		},
		{
			name:        "valid task creation with image",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			ownerID:     "user-1",
			imagePath:   "/uploads/images/test.jpg",
			wantErr:     false,
		},
		{
			name:        "empty title",
			title:       "",
			description: "Description",
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     true,
		},
		{
			name:        "empty owner id",
			title:       "Buy groceries",
			description: "Description",
			ownerID:     "",
			imagePath:   "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := useCase.Execute(context.Background(), tt.title, tt.description, tt.ownerID, tt.imagePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error = %v", err)
				return
			}

			if task == nil {
				t.Error("Execute() returned nil task")
				return
			}

			if task.Title != tt.title {
				t.Errorf("Task.Title = %v, want %v", task.Title, tt.title)
			}
			if task.Description != tt.description {
				t.Errorf("Task.Description = %v, want %v", task.Description, tt.description)
			}
			if task.OwnerID != tt.ownerID {
				t.Errorf("Task.OwnerID = %v, want %v", task.OwnerID, tt.ownerID)
			}
			if task.ImagePath != tt.imagePath {
				t.Errorf("Task.ImagePath = %v, want %v", task.ImagePath, tt.imagePath)
			}
			if task.Status != application.StatusPending {
				t.Errorf("Task.Status = %v, want %v", task.Status, application.StatusPending)
			}
		})
	}
}

// Mock repository
type mockTaskRepository struct {
	tasks map[string]*application.Task
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
	return m.tasks[id], nil
}

func (m *mockTaskRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*application.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) FindSharedWithUser(ctx context.Context, userID string) ([]*application.Task, error) {
	return nil, nil
}
