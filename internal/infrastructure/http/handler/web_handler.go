package handler

import (
	"net/http"

	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

// WebTaskHandler handles web requests (form data -> JSON)
type WebTaskHandler struct {
	createTask   *usecases.CreateTaskUseCase
	deleteTask   *usecases.DeleteTaskUseCase
	completeTask *usecases.CompleteTaskUseCase
}

// NewWebTaskHandler creates a new WebTaskHandler
func NewWebTaskHandler(
	createTask *usecases.CreateTaskUseCase,
	deleteTask *usecases.DeleteTaskUseCase,
	completeTask *usecases.CompleteTaskUseCase,
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
	html := `<div class="bg-white shadow rounded-lg p-6" id="task-` + task.ID + `">
		<div class="flex justify-between items-start">
			<div class="flex-1">
				<h3 class="text-lg font-semibold text-gray-900">` + task.Title + `</h3>
				<p class="text-gray-600 mt-1">` + task.Description + `</p>
				<div class="mt-2 flex items-center space-x-2">
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
						Pendente
					</span>
					<span class="text-sm text-gray-500">` + task.CreatedAt.Format("02/01/2006 15:04") + `</span>
				</div>
			</div>
			<div class="flex space-x-2 ml-4">
				<button hx-post="/web/tasks/` + task.ID + `/complete" hx-target="#task-` + task.ID + `" hx-swap="outerHTML"
						class="text-green-600 hover:text-green-800 font-medium">
					Concluir
				</button>
				<button hx-delete="/web/tasks/` + task.ID + `" hx-target="#task-` + task.ID + `" hx-swap="outerHTML"
						hx-confirm="Tem certeza que deseja excluir esta tarefa?"
						class="text-red-600 hover:text-red-800">
					Excluir
				</button>
			</div>
		</div>
	</div>`

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

	err := h.completeTask.Execute(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Return updated HTML fragment for HTMX with completed status
	w.Header().Set("Content-Type", "text/html")
	html := `<div class="bg-white shadow rounded-lg p-6" id="task-` + taskID + `">
		<div class="flex justify-between items-start">
			<div class="flex-1">
				<div class="flex items-center space-x-2">
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
						Concluída
					</span>
					<span class="text-sm text-gray-500">Tarefa concluída com sucesso!</span>
				</div>
			</div>
			<div class="flex space-x-2 ml-4">
				<button hx-delete="/web/tasks/` + taskID + `" hx-target="#task-` + taskID + `" hx-swap="outerHTML"
						hx-confirm="Tem certeza que deseja excluir esta tarefa?"
						class="text-red-600 hover:text-red-800">
					Excluir
				</button>
			</div>
		</div>
	</div>`

	w.Write([]byte(html))
}
