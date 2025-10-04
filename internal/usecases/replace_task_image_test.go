package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// Mock repositories for testing
type mockTaskRepositoryForReplaceImage struct {
	tasks map[string]*application.Task
}

func (m *mockTaskRepositoryForReplaceImage) Create(ctx context.Context, task *application.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForReplaceImage) Update(ctx context.Context, task *application.Task) error {
	if _, exists := m.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForReplaceImage) Delete(ctx context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepositoryForReplaceImage) FindByID(ctx context.Context, id string) (*application.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (m *mockTaskRepositoryForReplaceImage) FindByOwnerID(ctx context.Context, ownerID string) ([]*application.Task, error) {
	var tasks []*application.Task
	for _, task := range m.tasks {
		if task.OwnerID == ownerID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *mockTaskRepositoryForReplaceImage) FindSharedWithUser(ctx context.Context, userID string) ([]*application.Task, error) {
	return []*application.Task{}, nil
}

type mockTaskServiceForReplaceImage struct {
	canModify bool
}

func (m *mockTaskServiceForReplaceImage) CanUserAccessTask(ctx context.Context, taskID, userID string) (bool, error) {
	return true, nil
}

func (m *mockTaskServiceForReplaceImage) CanUserModifyTask(ctx context.Context, taskID, userID string) (bool, error) {
	return m.canModify, nil
}

func TestReplaceTaskImageUseCase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		taskID       string
		userID       string
		newImagePath string
		setupTask    func(*mockTaskRepositoryForReplaceImage)
		canModify    bool
		wantErr      bool
		errorMsg     string
		checkResult  func(*testing.T, string)
	}{
		{
			name:         "should replace image in pending task",
			taskID:       "task-1",
			userID:       "user-1",
			newImagePath: "/uploads/images/new.jpg",
			setupTask: func(repo *mockTaskRepositoryForReplaceImage) {
				task, _ := application.NewTask("task-1", "Test Task", "Description", application.StatusPending, "user-1", "/uploads/images/old.jpg")
				repo.tasks["task-1"] = task
			},
			canModify: true,
			wantErr:   false,
			checkResult: func(t *testing.T, oldImagePath string) {
				if oldImagePath != "/uploads/images/old.jpg" {
					t.Errorf("Execute() old image path = %v, want /uploads/images/old.jpg", oldImagePath)
				}
			},
		},
		{
			name:         "should replace image in in_progress task",
			taskID:       "task-2",
			userID:       "user-1",
			newImagePath: "/uploads/images/new.jpg",
			setupTask: func(repo *mockTaskRepositoryForReplaceImage) {
				task, _ := application.NewTask("task-2", "Test Task", "Description", application.StatusInProgress, "user-1", "/uploads/images/old.jpg")
				repo.tasks["task-2"] = task
			},
			canModify: true,
			wantErr:   false,
			checkResult: func(t *testing.T, oldImagePath string) {
				if oldImagePath != "/uploads/images/old.jpg" {
					t.Errorf("Execute() old image path = %v, want /uploads/images/old.jpg", oldImagePath)
				}
			},
		},
		{
			name:         "should fail if task not found",
			taskID:       "nonexistent",
			userID:       "user-1",
			newImagePath: "/uploads/images/new.jpg",
			setupTask: func(repo *mockTaskRepositoryForReplaceImage) {
				// No task setup
			},
			canModify: true,
			wantErr:   true,
			errorMsg:  "task not found",
		},
		{
			name:         "should fail if user cannot modify task",
			taskID:       "task-3",
			userID:       "user-2",
			newImagePath: "/uploads/images/new.jpg",
			setupTask: func(repo *mockTaskRepositoryForReplaceImage) {
				task, _ := application.NewTask("task-3", "Test Task", "Description", application.StatusPending, "user-1", "/uploads/images/old.jpg")
				repo.tasks["task-3"] = task
			},
			canModify: false,
			wantErr:   true,
			errorMsg:  "user does not have permission to modify this task",
		},
		{
			name:         "should fail if task is completed",
			taskID:       "task-4",
			userID:       "user-1",
			newImagePath: "/uploads/images/new.jpg",
			setupTask: func(repo *mockTaskRepositoryForReplaceImage) {
				task, _ := application.NewTask("task-4", "Test Task", "Description", application.StatusCompleted, "user-1", "/uploads/images/old.jpg")
				repo.tasks["task-4"] = task
			},
			canModify: true,
			wantErr:   true,
			errorMsg:  "cannot replace image in completed task",
		},
		{
			name:         "should fail if new image path is empty",
			taskID:       "task-5",
			userID:       "user-1",
			newImagePath: "",
			setupTask: func(repo *mockTaskRepositoryForReplaceImage) {
				task, _ := application.NewTask("task-5", "Test Task", "Description", application.StatusPending, "user-1", "/uploads/images/old.jpg")
				repo.tasks["task-5"] = task
			},
			canModify: true,
			wantErr:   true,
			errorMsg:  "new image path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockTaskRepositoryForReplaceImage{
				tasks: make(map[string]*application.Task),
			}
			tt.setupTask(mockRepo)

			mockService := &mockTaskServiceForReplaceImage{
				canModify: tt.canModify,
			}

			useCase := NewReplaceTaskImageUseCase(mockRepo, mockService)
			oldImagePath, err := useCase.Execute(context.Background(), tt.taskID, tt.userID, tt.newImagePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error but got nil")
					return
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

			// Verify task image was replaced
			task, _ := mockRepo.FindByID(context.Background(), tt.taskID)
			if task.ImagePath != tt.newImagePath {
				t.Errorf("Execute() task.ImagePath = %v, want %v", task.ImagePath, tt.newImagePath)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, oldImagePath)
			}
		})
	}
}
