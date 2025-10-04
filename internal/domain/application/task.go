package application

import (
	"errors"
	"time"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
)

// Task represents a todo task entity
type Task struct {
	ID          string
	Title       string
	Description string
	Status      TaskStatus
	OwnerID     string
	ImagePath   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask creates a new Task with validation
func NewTask(id, title, description string, status TaskStatus, ownerID, imagePath string) (*Task, error) {
	if id == "" {
		return nil, errors.New("task id cannot be empty")
	}

	if title == "" {
		return nil, errors.New("task title cannot be empty")
	}

	if len(title) > 200 {
		return nil, errors.New("task title cannot exceed 200 characters")
	}

	if len(description) > 1000 {
		return nil, errors.New("task description cannot exceed 1000 characters")
	}

	if ownerID == "" {
		return nil, errors.New("task owner id cannot be empty")
	}

	if len(imagePath) > 500 {
		return nil, errors.New("image path cannot exceed 500 characters")
	}

	if !isValidStatus(status) {
		return nil, errors.New("invalid task status")
	}

	now := time.Now()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
		OwnerID:     ownerID,
		ImagePath:   imagePath,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update updates task fields with validation
func (t *Task) Update(title, description string, status TaskStatus, imagePath string) error {
	if title == "" {
		return errors.New("task title cannot be empty")
	}

	if len(title) > 200 {
		return errors.New("task title cannot exceed 200 characters")
	}

	if len(description) > 1000 {
		return errors.New("task description cannot exceed 1000 characters")
	}

	if len(imagePath) > 500 {
		return errors.New("image path cannot exceed 500 characters")
	}

	if !isValidStatus(status) {
		return errors.New("invalid task status")
	}

	t.Title = title
	t.Description = description
	t.Status = status
	t.ImagePath = imagePath
	t.UpdatedAt = time.Now()

	return nil
}

// CompleteTask marks the task as completed
func (t *Task) CompleteTask() error {
	if t.Status == StatusCompleted {
		return errors.New("task is already completed")
	}

	t.Status = StatusCompleted
	t.UpdatedAt = time.Now()
	return nil
}

// RemoveImage removes the image from the task
func (t *Task) RemoveImage() error {
	if t.Status == StatusCompleted {
		return errors.New("cannot remove image from completed task")
	}

	if t.ImagePath == "" {
		return errors.New("task has no image to remove")
	}

	t.ImagePath = ""
	t.UpdatedAt = time.Now()
	return nil
}

// ReplaceImage replaces the current image with a new one
func (t *Task) ReplaceImage(newImagePath string) error {
	if t.Status == StatusCompleted {
		return errors.New("cannot replace image in completed task")
	}

	if newImagePath == "" {
		return errors.New("new image path cannot be empty")
	}

	if len(newImagePath) > 500 {
		return errors.New("image path cannot exceed 500 characters")
	}

	t.ImagePath = newImagePath
	t.UpdatedAt = time.Now()
	return nil
}

// isValidStatus checks if the status is valid
func isValidStatus(status TaskStatus) bool {
	return status == StatusPending || status == StatusInProgress || status == StatusCompleted
}
