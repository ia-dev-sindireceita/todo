package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	createTask      usecases.CreateTaskUseCaseInterface
	updateTask      usecases.UpdateTaskUseCaseInterface
	deleteTask      usecases.DeleteTaskUseCaseInterface
	getTask         usecases.GetTaskUseCaseInterface
	listTasks       usecases.ListTasksUseCaseInterface
	listSharedTasks usecases.ListSharedTasksUseCaseInterface
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(
	createTask usecases.CreateTaskUseCaseInterface,
	updateTask usecases.UpdateTaskUseCaseInterface,
	deleteTask usecases.DeleteTaskUseCaseInterface,
	getTask usecases.GetTaskUseCaseInterface,
	listTasks usecases.ListTasksUseCaseInterface,
	listSharedTasks usecases.ListSharedTasksUseCaseInterface,
) *TaskHandler {
	return &TaskHandler{
		createTask:      createTask,
		updateTask:      updateTask,
		deleteTask:      deleteTask,
		getTask:         getTask,
		listTasks:       listTasks,
		listSharedTasks: listSharedTasks,
	}
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// CreateTask handles POST /api/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(string)

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.createTask.Execute(r.Context(), req.Title, req.Description, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// ListTasks handles GET /api/tasks
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	tasks, err := h.listTasks.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// ListSharedTasks handles GET /api/tasks/shared
func (h *TaskHandler) ListSharedTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	tasks, err := h.listSharedTasks.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTask handles GET /api/tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("id")

	task, err := h.getTask.Execute(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTask handles PUT /api/tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("id")

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := application.TaskStatus(req.Status)
	err := h.updateTask.Execute(r.Context(), taskID, req.Title, req.Description, status, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTask handles DELETE /api/tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	taskID := r.PathValue("id")

	err := h.deleteTask.Execute(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
