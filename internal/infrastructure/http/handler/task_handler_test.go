package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// =============================================================================
// Mocks
// =============================================================================

type mockCreateTaskUseCase struct {
	executeFunc func(ctx context.Context, title, description, ownerID string) (*application.Task, error)
}

func (m *mockCreateTaskUseCase) Execute(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, title, description, ownerID)
	}
	return &application.Task{
		ID:          "task-123",
		Title:       title,
		Description: description,
		Status:      application.StatusPending,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

type mockGetTaskUseCase struct {
	executeFunc func(ctx context.Context, taskID, userID string) (*application.Task, error)
}

func (m *mockGetTaskUseCase) Execute(ctx context.Context, taskID, userID string) (*application.Task, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, taskID, userID)
	}
	return &application.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      application.StatusPending,
		OwnerID:     userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

type mockUpdateTaskUseCase struct {
	executeFunc func(ctx context.Context, taskID, title, description string, status application.TaskStatus, userID string) error
}

func (m *mockUpdateTaskUseCase) Execute(ctx context.Context, taskID, title, description string, status application.TaskStatus, userID string) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, taskID, title, description, status, userID)
	}
	return nil
}

type mockDeleteTaskUseCase struct {
	executeFunc func(ctx context.Context, taskID, userID string) error
}

func (m *mockDeleteTaskUseCase) Execute(ctx context.Context, taskID, userID string) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, taskID, userID)
	}
	return nil
}

type mockListTasksUseCase struct {
	executeFunc func(ctx context.Context, userID string) ([]*application.Task, error)
}

func (m *mockListTasksUseCase) Execute(ctx context.Context, userID string) ([]*application.Task, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, userID)
	}
	return []*application.Task{
		{
			ID:          "task-1",
			Title:       "Task 1",
			Description: "Description 1",
			Status:      application.StatusPending,
			OwnerID:     userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

type mockListSharedTasksUseCase struct {
	executeFunc func(ctx context.Context, userID string) ([]*application.Task, error)
}

func (m *mockListSharedTasksUseCase) Execute(ctx context.Context, userID string) ([]*application.Task, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, userID)
	}
	return []*application.Task{
		{
			ID:          "shared-task-1",
			Title:       "Shared Task 1",
			Description: "Shared Description 1",
			Status:      application.StatusPending,
			OwnerID:     "other-user",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil
}

// =============================================================================
// CreateTask Tests
// =============================================================================

func TestCreateTask_Success(t *testing.T) {
	mockCreate := &mockCreateTaskUseCase{
		executeFunc: func(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
			if title != "New Task" {
				t.Errorf("Expected title 'New Task', got %s", title)
			}
			if description != "Task description" {
				t.Errorf("Expected description 'Task description', got %s", description)
			}
			if ownerID != "user-123" {
				t.Errorf("Expected ownerID 'user-123', got %s", ownerID)
			}
			return &application.Task{
				ID:          "task-456",
				Title:       title,
				Description: description,
				Status:      application.StatusPending,
				OwnerID:     ownerID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		},
	}

	handler := NewTaskHandler(mockCreate, nil, nil, nil, nil, nil)

	reqBody := CreateTaskRequest{
		Title:       "New Task",
		Description: "Task description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response application.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != "task-456" {
		t.Errorf("Expected task ID 'task-456', got %s", response.ID)
	}

	if response.Title != "New Task" {
		t.Errorf("Expected title 'New Task', got %s", response.Title)
	}
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	handler := NewTaskHandler(&mockCreateTaskUseCase{}, nil, nil, nil, nil, nil)

	req := httptest.NewRequest("POST", "/api/tasks", strings.NewReader("invalid-json"))
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Invalid request body") {
		t.Errorf("Expected error message about invalid request body, got: %s", body)
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	mockCreate := &mockCreateTaskUseCase{
		executeFunc: func(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
			return nil, errors.New("task title cannot be empty")
		},
	}

	handler := NewTaskHandler(mockCreate, nil, nil, nil, nil, nil)

	reqBody := CreateTaskRequest{
		Title:       "",
		Description: "Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "task title cannot be empty") {
		t.Errorf("Expected validation error, got: %s", responseBody)
	}
}

func TestCreateTask_TitleTooLong(t *testing.T) {
	mockCreate := &mockCreateTaskUseCase{
		executeFunc: func(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
			if len(title) > 200 {
				return nil, errors.New("task title cannot exceed 200 characters")
			}
			return nil, nil
		},
	}

	handler := NewTaskHandler(mockCreate, nil, nil, nil, nil, nil)

	longTitle := strings.Repeat("a", 201)
	reqBody := CreateTaskRequest{
		Title:       longTitle,
		Description: "Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// =============================================================================
// GetTask Tests
// =============================================================================

func TestGetTask_Success(t *testing.T) {
	mockGet := &mockGetTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			if taskID != "task-123" {
				t.Errorf("Expected taskID 'task-123', got %s", taskID)
			}
			if userID != "user-123" {
				t.Errorf("Expected userID 'user-123', got %s", userID)
			}
			return &application.Task{
				ID:          taskID,
				Title:       "Test Task",
				Description: "Test Description",
				Status:      application.StatusPending,
				OwnerID:     userID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		},
	}

	handler := NewTaskHandler(nil, nil, nil, mockGet, nil, nil)

	req := httptest.NewRequest("GET", "/api/tasks/task-123", nil)
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response application.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != "task-123" {
		t.Errorf("Expected task ID 'task-123', got %s", response.ID)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	mockGet := &mockGetTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			return nil, errors.New("task not found")
		},
	}

	handler := NewTaskHandler(nil, nil, nil, mockGet, nil, nil)

	req := httptest.NewRequest("GET", "/api/tasks/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestGetTask_NoPermission(t *testing.T) {
	mockGet := &mockGetTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			return nil, errors.New("user does not have permission to access this task")
		},
	}

	handler := NewTaskHandler(nil, nil, nil, mockGet, nil, nil)

	req := httptest.NewRequest("GET", "/api/tasks/task-123", nil)
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "other-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

// =============================================================================
// UpdateTask Tests
// =============================================================================

func TestUpdateTask_Success(t *testing.T) {
	mockUpdate := &mockUpdateTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, title, description string, status application.TaskStatus, userID string) error {
			if taskID != "task-123" {
				t.Errorf("Expected taskID 'task-123', got %s", taskID)
			}
			if title != "Updated Task" {
				t.Errorf("Expected title 'Updated Task', got %s", title)
			}
			if status != application.StatusInProgress {
				t.Errorf("Expected status in_progress, got %s", status)
			}
			return nil
		},
	}

	handler := NewTaskHandler(nil, mockUpdate, nil, nil, nil, nil)

	reqBody := UpdateTaskRequest{
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      "in_progress",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/api/tasks/task-123", bytes.NewReader(body))
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UpdateTask(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestUpdateTask_InvalidJSON(t *testing.T) {
	handler := NewTaskHandler(nil, &mockUpdateTaskUseCase{}, nil, nil, nil, nil)

	req := httptest.NewRequest("PUT", "/api/tasks/task-123", strings.NewReader("invalid-json"))
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTask_InvalidStatus(t *testing.T) {
	mockUpdate := &mockUpdateTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, title, description string, status application.TaskStatus, userID string) error {
			return errors.New("invalid task status")
		},
	}

	handler := NewTaskHandler(nil, mockUpdate, nil, nil, nil, nil)

	reqBody := UpdateTaskRequest{
		Title:       "Task",
		Description: "Description",
		Status:      "invalid_status",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/api/tasks/task-123", bytes.NewReader(body))
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTask_NoPermission(t *testing.T) {
	mockUpdate := &mockUpdateTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, title, description string, status application.TaskStatus, userID string) error {
			return errors.New("user does not have permission to modify this task")
		},
	}

	handler := NewTaskHandler(nil, mockUpdate, nil, nil, nil, nil)

	reqBody := UpdateTaskRequest{
		Title:       "Task",
		Description: "Description",
		Status:      "pending",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/api/tasks/task-123", bytes.NewReader(body))
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "other-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// =============================================================================
// DeleteTask Tests
// =============================================================================

func TestDeleteTask_Success(t *testing.T) {
	mockDelete := &mockDeleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) error {
			if taskID != "task-123" {
				t.Errorf("Expected taskID 'task-123', got %s", taskID)
			}
			if userID != "user-123" {
				t.Errorf("Expected userID 'user-123', got %s", userID)
			}
			return nil
		},
	}

	handler := NewTaskHandler(nil, nil, mockDelete, nil, nil, nil)

	req := httptest.NewRequest("DELETE", "/api/tasks/task-123", nil)
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	mockDelete := &mockDeleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) error {
			return errors.New("task not found")
		},
	}

	handler := NewTaskHandler(nil, nil, mockDelete, nil, nil, nil)

	req := httptest.NewRequest("DELETE", "/api/tasks/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestDeleteTask_NoPermission(t *testing.T) {
	mockDelete := &mockDeleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) error {
			return errors.New("user does not have permission to delete this task")
		},
	}

	handler := NewTaskHandler(nil, nil, mockDelete, nil, nil, nil)

	req := httptest.NewRequest("DELETE", "/api/tasks/task-123", nil)
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "other-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

// =============================================================================
// ListTasks Tests
// =============================================================================

func TestListTasks_Success(t *testing.T) {
	mockList := &mockListTasksUseCase{
		executeFunc: func(ctx context.Context, userID string) ([]*application.Task, error) {
			if userID != "user-123" {
				t.Errorf("Expected userID 'user-123', got %s", userID)
			}
			return []*application.Task{
				{
					ID:          "task-1",
					Title:       "Task 1",
					Description: "Description 1",
					Status:      application.StatusPending,
					OwnerID:     userID,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          "task-2",
					Title:       "Task 2",
					Description: "Description 2",
					Status:      application.StatusCompleted,
					OwnerID:     userID,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			}, nil
		},
	}

	handler := NewTaskHandler(nil, nil, nil, nil, mockList, nil)

	req := httptest.NewRequest("GET", "/api/tasks", nil)
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response []*application.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(response))
	}

	if response[0].ID != "task-1" {
		t.Errorf("Expected first task ID 'task-1', got %s", response[0].ID)
	}
}

func TestListTasks_Empty(t *testing.T) {
	mockList := &mockListTasksUseCase{
		executeFunc: func(ctx context.Context, userID string) ([]*application.Task, error) {
			return []*application.Task{}, nil
		},
	}

	handler := NewTaskHandler(nil, nil, nil, nil, mockList, nil)

	req := httptest.NewRequest("GET", "/api/tasks", nil)
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*application.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(response))
	}
}

func TestListTasks_Error(t *testing.T) {
	mockList := &mockListTasksUseCase{
		executeFunc: func(ctx context.Context, userID string) ([]*application.Task, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewTaskHandler(nil, nil, nil, nil, mockList, nil)

	req := httptest.NewRequest("GET", "/api/tasks", nil)
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListTasks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// =============================================================================
// ListSharedTasks Tests
// =============================================================================

func TestListSharedTasks_Success(t *testing.T) {
	mockListShared := &mockListSharedTasksUseCase{
		executeFunc: func(ctx context.Context, userID string) ([]*application.Task, error) {
			if userID != "user-123" {
				t.Errorf("Expected userID 'user-123', got %s", userID)
			}
			return []*application.Task{
				{
					ID:          "shared-task-1",
					Title:       "Shared Task 1",
					Description: "Shared Description 1",
					Status:      application.StatusPending,
					OwnerID:     "other-user",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			}, nil
		},
	}

	handler := NewTaskHandler(nil, nil, nil, nil, nil, mockListShared)

	req := httptest.NewRequest("GET", "/api/tasks/shared", nil)
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListSharedTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response []*application.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 task, got %d", len(response))
	}

	if response[0].OwnerID == "user-123" {
		t.Error("Expected shared task from different owner")
	}
}

func TestListSharedTasks_Empty(t *testing.T) {
	mockListShared := &mockListSharedTasksUseCase{
		executeFunc: func(ctx context.Context, userID string) ([]*application.Task, error) {
			return []*application.Task{}, nil
		},
	}

	handler := NewTaskHandler(nil, nil, nil, nil, nil, mockListShared)

	req := httptest.NewRequest("GET", "/api/tasks/shared", nil)
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListSharedTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*application.Task
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(response))
	}
}

func TestListSharedTasks_Error(t *testing.T) {
	mockListShared := &mockListSharedTasksUseCase{
		executeFunc: func(ctx context.Context, userID string) ([]*application.Task, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewTaskHandler(nil, nil, nil, nil, nil, mockListShared)

	req := httptest.NewRequest("GET", "/api/tasks/shared", nil)
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListSharedTasks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
