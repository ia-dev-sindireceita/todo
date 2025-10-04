package handler

import (
	"net/http"

	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

// WebTaskHandler handles web requests (form data -> JSON)
type WebTaskHandler struct {
	createTask   usecases.CreateTaskUseCaseInterface
	deleteTask   usecases.DeleteTaskUseCaseInterface
	completeTask usecases.CompleteTaskUseCaseInterface
}

// NewWebTaskHandler creates a new WebTaskHandler
func NewWebTaskHandler(
	createTask usecases.CreateTaskUseCaseInterface,
	deleteTask usecases.DeleteTaskUseCaseInterface,
	completeTask usecases.CompleteTaskUseCaseInterface,
) *WebTaskHandler {
	return &WebTaskHandler{
		createTask:   createTask,
		deleteTask:   deleteTask,
		completeTask: completeTask,
	}
}

// CreateTask handles web form submission
func (h *WebTaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	// Create task
	task, err := h.createTask.Execute(r.Context(), title, description, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return HTML fragment for HTMX
	w.Header().Set("Content-Type", "text/html")
	html, err := renderTaskCard(task, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(html))
}

// DeleteTask handles task deletion
func (h *WebTaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := r.PathValue("id")

	err := h.deleteTask.Execute(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Return empty response for HTMX to swap out the element
	w.WriteHeader(http.StatusOK)
}

// CompleteTask handles task completion
func (h *WebTaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := r.PathValue("id")

	task, err := h.completeTask.Execute(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Return updated HTML fragment for HTMX with completed status
	w.Header().Set("Content-Type", "text/html")
	html, err := renderCompletedTask(task, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(html))
}
