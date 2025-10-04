package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// =============================================================================
// WebCreateTask Tests
// =============================================================================

func TestWebCreateTask_Success(t *testing.T) {
	mockCreate := &mockCreateTaskUseCase{
		executeFunc: func(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
			if title != "New Web Task" {
				t.Errorf("Expected title 'New Web Task', got %s", title)
			}
			if description != "Web task description" {
				t.Errorf("Expected description 'Web task description', got %s", description)
			}
			if ownerID != "user-123" {
				t.Errorf("Expected ownerID 'user-123', got %s", ownerID)
			}
			return &application.Task{
				ID:          "web-task-456",
				Title:       title,
				Description: description,
				Status:      application.StatusPending,
				OwnerID:     ownerID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		},
	}

	handler := NewWebTaskHandler(mockCreate, nil, nil)

	formData := url.Values{}
	formData.Set("title", "New Web Task")
	formData.Set("description", "Web task description")

	req := httptest.NewRequest("POST", "/web/tasks", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()
	if !strings.Contains(body, "New Web Task") {
		t.Errorf("Expected HTML fragment to contain task title, got: %s", body)
	}

	if !strings.Contains(body, "Web task description") {
		t.Errorf("Expected HTML fragment to contain task description, got: %s", body)
	}

	if !strings.Contains(body, "web-task-456") {
		t.Errorf("Expected HTML fragment to contain task ID, got: %s", body)
	}

	// Verify HTMX attributes
	if !strings.Contains(body, "hx-post") {
		t.Error("Expected HTML fragment to contain HTMX hx-post attribute")
	}

	if !strings.Contains(body, "hx-delete") {
		t.Error("Expected HTML fragment to contain HTMX hx-delete attribute")
	}

	// Verify Tailwind CSS classes
	if !strings.Contains(body, "bg-white") {
		t.Error("Expected HTML fragment to contain Tailwind CSS classes")
	}

	// Verify status badge
	if !strings.Contains(body, "Pendente") {
		t.Error("Expected HTML fragment to contain 'Pendente' status")
	}

	if !strings.Contains(body, "bg-yellow-100") {
		t.Error("Expected HTML fragment to contain yellow status badge")
	}

	// Verify ownership badge for own task
	if !strings.Contains(body, "Própria") {
		t.Error("Expected HTML fragment to contain 'Própria' ownership badge")
	}

	if !strings.Contains(body, "bg-blue-100") {
		t.Error("Expected HTML fragment to contain blue ownership badge for own task")
	}
}

func TestWebCreateTask_SharedTaskIndicator(t *testing.T) {
	mockCreate := &mockCreateTaskUseCase{
		executeFunc: func(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
			// Simula que outro usuário criou a tarefa
			return &application.Task{
				ID:          "shared-task-789",
				Title:       "Shared Task",
				Description: "Task shared with me",
				Status:      application.StatusPending,
				OwnerID:     "other-user-456", // Diferente do usuário atual
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		},
	}

	handler := NewWebTaskHandler(mockCreate, nil, nil)

	formData := url.Values{}
	formData.Set("title", "Shared Task")
	formData.Set("description", "Task shared with me")

	req := httptest.NewRequest("POST", "/web/tasks", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	body := w.Body.String()

	// Verify shared task badge
	if !strings.Contains(body, "Compartilhada") {
		t.Error("Expected HTML fragment to contain 'Compartilhada' ownership badge")
	}

	if !strings.Contains(body, "bg-purple-100") {
		t.Error("Expected HTML fragment to contain purple ownership badge for shared task")
	}
}

func TestWebCreateTask_Unauthorized(t *testing.T) {
	handler := NewWebTaskHandler(&mockCreateTaskUseCase{}, nil, nil)

	formData := url.Values{}
	formData.Set("title", "Task")
	formData.Set("description", "Description")

	req := httptest.NewRequest("POST", "/web/tasks", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// No userID in context

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Unauthorized") {
		t.Errorf("Expected error message 'Unauthorized', got: %s", body)
	}
}

func TestWebCreateTask_InvalidForm(t *testing.T) {
	handler := NewWebTaskHandler(&mockCreateTaskUseCase{}, nil, nil)

	req := httptest.NewRequest("POST", "/web/tasks", strings.NewReader("%invalid%form%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Invalid form data") {
		t.Errorf("Expected error message about invalid form data, got: %s", body)
	}
}

func TestWebCreateTask_ValidationError(t *testing.T) {
	mockCreate := &mockCreateTaskUseCase{
		executeFunc: func(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
			return nil, errors.New("task title cannot be empty")
		},
	}

	handler := NewWebTaskHandler(mockCreate, nil, nil)

	formData := url.Values{}
	formData.Set("title", "")
	formData.Set("description", "Description")

	req := httptest.NewRequest("POST", "/web/tasks", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "task title cannot be empty") {
		t.Errorf("Expected validation error message, got: %s", body)
	}
}

func TestWebCreateTask_HTMLEscaping(t *testing.T) {
	mockCreate := &mockCreateTaskUseCase{
		executeFunc: func(ctx context.Context, title, description, ownerID string) (*application.Task, error) {
			return &application.Task{
				ID:          "task-xss",
				Title:       title,
				Description: description,
				Status:      application.StatusPending,
				OwnerID:     ownerID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		},
	}

	handler := NewWebTaskHandler(mockCreate, nil, nil)

	// Test with potentially malicious input
	formData := url.Values{}
	formData.Set("title", "<script>alert('xss')</script>")
	formData.Set("description", "<img src=x onerror=alert('xss')>")

	req := httptest.NewRequest("POST", "/web/tasks", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CreateTask(w, req)

	body := w.Body.String()

	// Verify that HTML is properly escaped
	if strings.Contains(body, "<script>") {
		t.Error("HTML content is not escaped - XSS vulnerability detected")
	}

	// Verify escaped content is present
	if !strings.Contains(body, "&lt;script&gt;") {
		t.Error("Expected script tags to be HTML escaped as &lt;script&gt;")
	}

	if !strings.Contains(body, "&lt;img") {
		t.Error("Expected img tag to be HTML escaped")
	}

	// Verify the content is still displayed (but safely)
	if !strings.Contains(body, "alert") {
		t.Error("Expected escaped content to still contain 'alert' text")
	}
}

// =============================================================================
// WebDeleteTask Tests
// =============================================================================

func TestWebDeleteTask_Success(t *testing.T) {
	mockDelete := &mockDeleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) error {
			if taskID != "task-to-delete" {
				t.Errorf("Expected taskID 'task-to-delete', got %s", taskID)
			}
			if userID != "user-123" {
				t.Errorf("Expected userID 'user-123', got %s", userID)
			}
			return nil
		},
	}

	handler := NewWebTaskHandler(nil, mockDelete, nil)

	req := httptest.NewRequest("DELETE", "/web/tasks/task-to-delete", nil)
	req.SetPathValue("id", "task-to-delete")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// For HTMX swap, empty response is expected
	body := w.Body.String()
	if body != "" {
		t.Errorf("Expected empty response body, got: %s", body)
	}
}

func TestWebDeleteTask_Unauthorized(t *testing.T) {
	handler := NewWebTaskHandler(nil, &mockDeleteTaskUseCase{}, nil)

	req := httptest.NewRequest("DELETE", "/web/tasks/task-123", nil)
	req.SetPathValue("id", "task-123")
	// No userID in context

	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestWebDeleteTask_NotFound(t *testing.T) {
	mockDelete := &mockDeleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) error {
			return errors.New("task not found")
		},
	}

	handler := NewWebTaskHandler(nil, mockDelete, nil)

	req := httptest.NewRequest("DELETE", "/web/tasks/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestWebDeleteTask_NoPermission(t *testing.T) {
	mockDelete := &mockDeleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) error {
			return errors.New("user does not have permission to delete this task")
		},
	}

	handler := NewWebTaskHandler(nil, mockDelete, nil)

	req := httptest.NewRequest("DELETE", "/web/tasks/task-123", nil)
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "other-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "permission") {
		t.Errorf("Expected permission error message, got: %s", body)
	}
}

// =============================================================================
// WebCompleteTask Tests
// =============================================================================

func TestWebCompleteTask_Success(t *testing.T) {
	mockComplete := &mockCompleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			if taskID != "task-to-complete" {
				t.Errorf("Expected taskID 'task-to-complete', got %s", taskID)
			}
			if userID != "user-123" {
				t.Errorf("Expected userID 'user-123', got %s", userID)
			}
			return &application.Task{
				ID:        "task-to-complete",
				Title:     "Test Task",
				OwnerID:   "user-123",
				Status:    application.StatusCompleted,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewWebTaskHandler(nil, nil, mockComplete)

	req := httptest.NewRequest("POST", "/web/tasks/task-to-complete/complete", nil)
	req.SetPathValue("id", "task-to-complete")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CompleteTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()

	// Verify completed status
	if !strings.Contains(body, "Concluída") {
		t.Error("Expected HTML fragment to contain 'Concluída' status")
	}

	if !strings.Contains(body, "bg-green-100") {
		t.Error("Expected HTML fragment to contain green status badge")
	}

	if !strings.Contains(body, "task-to-complete") {
		t.Errorf("Expected HTML fragment to contain task ID, got: %s", body)
	}

	// Verify success message
	if !strings.Contains(body, "Tarefa concluída com sucesso") {
		t.Error("Expected success message in HTML fragment")
	}

	// Verify only delete button remains (no complete button)
	if strings.Contains(body, "Concluir") {
		t.Error("Expected HTML fragment to NOT contain 'Concluir' button")
	}

	if !strings.Contains(body, "hx-delete") {
		t.Error("Expected HTML fragment to contain delete button")
	}

	// Verify Tailwind CSS classes
	if !strings.Contains(body, "bg-white") {
		t.Error("Expected HTML fragment to contain Tailwind CSS classes")
	}

	// Verify ownership badge for completed own task
	if !strings.Contains(body, "Própria") {
		t.Error("Expected HTML fragment to contain 'Própria' ownership badge")
	}

	if !strings.Contains(body, "bg-blue-100") {
		t.Error("Expected HTML fragment to contain blue ownership badge for own task")
	}
}

func TestWebCompleteTask_SharedTaskIndicator(t *testing.T) {
	mockComplete := &mockCompleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			return &application.Task{
				ID:        "shared-task-999",
				Title:     "Shared Task",
				OwnerID:   "other-user-456", // Different from current user
				Status:    application.StatusCompleted,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewWebTaskHandler(nil, nil, mockComplete)

	req := httptest.NewRequest("POST", "/web/tasks/shared-task-999/complete", nil)
	req.SetPathValue("id", "shared-task-999")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CompleteTask(w, req)

	body := w.Body.String()

	// Verify shared ownership badge in completed task
	if !strings.Contains(body, "Compartilhada") {
		t.Error("Expected completed shared task to contain 'Compartilhada' ownership badge")
	}

	if !strings.Contains(body, "bg-purple-100") {
		t.Error("Expected completed shared task to contain purple ownership badge")
	}
}

func TestWebCompleteTask_Unauthorized(t *testing.T) {
	handler := NewWebTaskHandler(nil, nil, &mockCompleteTaskUseCase{})

	req := httptest.NewRequest("POST", "/web/tasks/task-123/complete", nil)
	req.SetPathValue("id", "task-123")
	// No userID in context

	w := httptest.NewRecorder()
	handler.CompleteTask(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestWebCompleteTask_NotFound(t *testing.T) {
	mockComplete := &mockCompleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			return nil, errors.New("task not found")
		},
	}

	handler := NewWebTaskHandler(nil, nil, mockComplete)

	req := httptest.NewRequest("POST", "/web/tasks/nonexistent/complete", nil)
	req.SetPathValue("id", "nonexistent")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CompleteTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestWebCompleteTask_NoPermission(t *testing.T) {
	mockComplete := &mockCompleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			return nil, errors.New("user does not have permission to modify this task")
		},
	}

	handler := NewWebTaskHandler(nil, nil, mockComplete)

	req := httptest.NewRequest("POST", "/web/tasks/task-123/complete", nil)
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "other-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CompleteTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestWebCompleteTask_AlreadyCompleted(t *testing.T) {
	mockComplete := &mockCompleteTaskUseCase{
		executeFunc: func(ctx context.Context, taskID, userID string) (*application.Task, error) {
			return nil, errors.New("task is already completed")
		},
	}

	handler := NewWebTaskHandler(nil, nil, mockComplete)

	req := httptest.NewRequest("POST", "/web/tasks/task-123/complete", nil)
	req.SetPathValue("id", "task-123")
	ctx := context.WithValue(req.Context(), "userID", "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.CompleteTask(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "already completed") {
		t.Errorf("Expected 'already completed' error message, got: %s", body)
	}
}

// Mock for CompleteTaskUseCase (needed for web handler tests)
type mockCompleteTaskUseCase struct {
	executeFunc func(ctx context.Context, taskID, userID string) (*application.Task, error)
}

func (m *mockCompleteTaskUseCase) Execute(ctx context.Context, taskID, userID string) (*application.Task, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, taskID, userID)
	}
	return nil, nil
}
