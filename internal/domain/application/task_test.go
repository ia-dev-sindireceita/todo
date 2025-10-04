package application

import (
	"strings"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		title       string
		description string
		status      TaskStatus
		ownerID     string
		imagePath   string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid task",
			id:          "task-1",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			status:      StatusPending,
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     false,
		},
		{
			name:        "valid task with image",
			id:          "task-2",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			status:      StatusPending,
			ownerID:     "user-1",
			imagePath:   "/uploads/images/abc123.jpg",
			wantErr:     false,
		},
		{
			name:        "image path too long",
			id:          "task-3",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			status:      StatusPending,
			ownerID:     "user-1",
			imagePath:   strings.Repeat("a", 501),
			wantErr:     true,
			errMsg:      "image path cannot exceed 500 characters",
		},
		{
			name:        "empty id",
			id:          "",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			status:      StatusPending,
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     true,
			errMsg:      "task id cannot be empty",
		},
		{
			name:        "empty title",
			id:          "task-1",
			title:       "",
			description: "Milk, bread, eggs",
			status:      StatusPending,
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     true,
			errMsg:      "task title cannot be empty",
		},
		{
			name:        "title too long",
			id:          "task-1",
			title:       string(make([]byte, 201)),
			description: "Milk, bread, eggs",
			status:      StatusPending,
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     true,
			errMsg:      "task title cannot exceed 200 characters",
		},
		{
			name:        "description too long",
			id:          "task-1",
			title:       "Buy groceries",
			description: string(make([]byte, 1001)),
			status:      StatusPending,
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     true,
			errMsg:      "task description cannot exceed 1000 characters",
		},
		{
			name:        "empty owner id",
			id:          "task-1",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			status:      StatusPending,
			ownerID:     "",
			imagePath:   "",
			wantErr:     true,
			errMsg:      "task owner id cannot be empty",
		},
		{
			name:        "invalid status",
			id:          "task-1",
			title:       "Buy groceries",
			description: "Milk, bread, eggs",
			status:      TaskStatus("invalid"),
			ownerID:     "user-1",
			imagePath:   "",
			wantErr:     true,
			errMsg:      "invalid task status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.id, tt.title, tt.description, tt.status, tt.ownerID, tt.imagePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewTask() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("NewTask() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTask() unexpected error = %v", err)
				return
			}

			if task.ID != tt.id {
				t.Errorf("Task.ID = %v, want %v", task.ID, tt.id)
			}
			if task.Title != tt.title {
				t.Errorf("Task.Title = %v, want %v", task.Title, tt.title)
			}
			if task.Description != tt.description {
				t.Errorf("Task.Description = %v, want %v", task.Description, tt.description)
			}
			if task.Status != tt.status {
				t.Errorf("Task.Status = %v, want %v", task.Status, tt.status)
			}
			if task.OwnerID != tt.ownerID {
				t.Errorf("Task.OwnerID = %v, want %v", task.OwnerID, tt.ownerID)
			}
			if task.ImagePath != tt.imagePath {
				t.Errorf("Task.ImagePath = %v, want %v", task.ImagePath, tt.imagePath)
			}
			if task.CreatedAt.IsZero() {
				t.Error("Task.CreatedAt should not be zero")
			}
			if task.UpdatedAt.IsZero() {
				t.Error("Task.UpdatedAt should not be zero")
			}
		})
	}
}

func TestTask_Update(t *testing.T) {
	task, err := NewTask("task-1", "Original title", "Original description", StatusPending, "user-1", "")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	originalUpdatedAt := task.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name        string
		title       string
		description string
		status      TaskStatus
		imagePath   string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid update",
			title:       "Updated title",
			description: "Updated description",
			status:      StatusInProgress,
			imagePath:   "",
			wantErr:     false,
		},
		{
			name:        "valid update with image",
			title:       "Updated title",
			description: "Updated description",
			status:      StatusInProgress,
			imagePath:   "/uploads/images/new.jpg",
			wantErr:     false,
		},
		{
			name:        "image path too long",
			title:       "Updated title",
			description: "Updated description",
			status:      StatusInProgress,
			imagePath:   strings.Repeat("a", 501),
			wantErr:     true,
			errMsg:      "image path cannot exceed 500 characters",
		},
		{
			name:        "empty title",
			title:       "",
			description: "Updated description",
			status:      StatusInProgress,
			imagePath:   "",
			wantErr:     true,
			errMsg:      "task title cannot be empty",
		},
		{
			name:        "invalid status",
			title:       "Updated title",
			description: "Updated description",
			status:      TaskStatus("invalid"),
			imagePath:   "",
			wantErr:     true,
			errMsg:      "invalid task status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := task.Update(tt.title, tt.description, tt.status, tt.imagePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Task.Update() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Task.Update() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Task.Update() unexpected error = %v", err)
				return
			}

			if task.Title != tt.title {
				t.Errorf("Task.Title = %v, want %v", task.Title, tt.title)
			}
			if task.Description != tt.description {
				t.Errorf("Task.Description = %v, want %v", task.Description, tt.description)
			}
			if task.Status != tt.status {
				t.Errorf("Task.Status = %v, want %v", task.Status, tt.status)
			}
			if task.ImagePath != tt.imagePath {
				t.Errorf("Task.ImagePath = %v, want %v", task.ImagePath, tt.imagePath)
			}
			if !task.UpdatedAt.After(originalUpdatedAt) {
				t.Error("Task.UpdatedAt should be updated")
			}
		})
	}
}

func TestTask_CompleteTask(t *testing.T) {
	tests := []struct {
		name    string
		status  TaskStatus
		wantErr bool
		errMsg  string
	}{
		{
			name:    "should complete pending task",
			status:  StatusPending,
			wantErr: false,
		},
		{
			name:    "should complete in_progress task",
			status:  StatusInProgress,
			wantErr: false,
		},
		{
			name:    "should fail if task already completed",
			status:  StatusCompleted,
			wantErr: true,
			errMsg:  "task is already completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("task-1", "Test Task", "Description", tt.status, "user-1", "")
			oldUpdatedAt := task.UpdatedAt

			err := task.CompleteTask()

			if tt.wantErr {
				if err == nil {
					t.Errorf("CompleteTask() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("CompleteTask() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("CompleteTask() unexpected error: %v", err)
			}

			if task.Status != StatusCompleted {
				t.Errorf("CompleteTask() status = %v, want %v", task.Status, StatusCompleted)
			}

			if !task.UpdatedAt.After(oldUpdatedAt) {
				t.Errorf("CompleteTask() did not update UpdatedAt")
			}
		})
	}
}

func TestTask_RemoveImage(t *testing.T) {
	tests := []struct {
		name      string
		status    TaskStatus
		imagePath string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "should remove image from pending task",
			status:    StatusPending,
			imagePath: "/uploads/images/test.jpg",
			wantErr:   false,
		},
		{
			name:      "should remove image from in_progress task",
			status:    StatusInProgress,
			imagePath: "/uploads/images/test.jpg",
			wantErr:   false,
		},
		{
			name:      "should fail to remove image from completed task",
			status:    StatusCompleted,
			imagePath: "/uploads/images/test.jpg",
			wantErr:   true,
			errMsg:    "cannot remove image from completed task",
		},
		{
			name:      "should fail when no image exists",
			status:    StatusPending,
			imagePath: "",
			wantErr:   true,
			errMsg:    "task has no image to remove",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("task-1", "Test Task", "Description", tt.status, "user-1", tt.imagePath)
			oldUpdatedAt := task.UpdatedAt
			time.Sleep(1 * time.Millisecond)

			err := task.RemoveImage()

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoveImage() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("RemoveImage() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("RemoveImage() unexpected error: %v", err)
			}

			if task.ImagePath != "" {
				t.Errorf("RemoveImage() ImagePath = %v, want empty string", task.ImagePath)
			}

			if !task.UpdatedAt.After(oldUpdatedAt) {
				t.Errorf("RemoveImage() did not update UpdatedAt")
			}
		})
	}
}

func TestTask_ReplaceImage(t *testing.T) {
	tests := []struct {
		name         string
		status       TaskStatus
		oldImagePath string
		newImagePath string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "should replace image in pending task",
			status:       StatusPending,
			oldImagePath: "/uploads/images/old.jpg",
			newImagePath: "/uploads/images/new.jpg",
			wantErr:      false,
		},
		{
			name:         "should replace image in in_progress task",
			status:       StatusInProgress,
			oldImagePath: "/uploads/images/old.jpg",
			newImagePath: "/uploads/images/new.jpg",
			wantErr:      false,
		},
		{
			name:         "should fail to replace image in completed task",
			status:       StatusCompleted,
			oldImagePath: "/uploads/images/old.jpg",
			newImagePath: "/uploads/images/new.jpg",
			wantErr:      true,
			errMsg:       "cannot replace image in completed task",
		},
		{
			name:         "should fail when new image path is empty",
			status:       StatusPending,
			oldImagePath: "/uploads/images/old.jpg",
			newImagePath: "",
			wantErr:      true,
			errMsg:       "new image path cannot be empty",
		},
		{
			name:         "should fail when new image path is too long",
			status:       StatusPending,
			oldImagePath: "/uploads/images/old.jpg",
			newImagePath: strings.Repeat("a", 501),
			wantErr:      true,
			errMsg:       "image path cannot exceed 500 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, _ := NewTask("task-1", "Test Task", "Description", tt.status, "user-1", tt.oldImagePath)
			oldUpdatedAt := task.UpdatedAt
			time.Sleep(1 * time.Millisecond)

			err := task.ReplaceImage(tt.newImagePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ReplaceImage() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ReplaceImage() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ReplaceImage() unexpected error: %v", err)
			}

			if task.ImagePath != tt.newImagePath {
				t.Errorf("ReplaceImage() ImagePath = %v, want %v", task.ImagePath, tt.newImagePath)
			}

			if !task.UpdatedAt.After(oldUpdatedAt) {
				t.Errorf("ReplaceImage() did not update UpdatedAt")
			}
		})
	}
}

