package handler

import (
	"net/http"

	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

// WebTaskHandler handles web requests (form data -> JSON)
type WebTaskHandler struct {
	createTask       usecases.CreateTaskUseCaseInterface
	deleteTask       usecases.DeleteTaskUseCaseInterface
	completeTask     usecases.CompleteTaskUseCaseInterface
	shareTask        usecases.ShareTaskUseCaseInterface
	deleteTaskImage  usecases.DeleteTaskImageUseCaseInterface
	replaceTaskImage usecases.ReplaceTaskImageUseCaseInterface
}

// NewWebTaskHandler creates a new WebTaskHandler
func NewWebTaskHandler(
	createTask usecases.CreateTaskUseCaseInterface,
	deleteTask usecases.DeleteTaskUseCaseInterface,
	completeTask usecases.CompleteTaskUseCaseInterface,
	shareTask usecases.ShareTaskUseCaseInterface,
	deleteTaskImage usecases.DeleteTaskImageUseCaseInterface,
	replaceTaskImage usecases.ReplaceTaskImageUseCaseInterface,
) *WebTaskHandler {
	return &WebTaskHandler{
		createTask:       createTask,
		deleteTask:       deleteTask,
		completeTask:     completeTask,
		shareTask:        shareTask,
		deleteTaskImage:  deleteTaskImage,
		replaceTaskImage: replaceTaskImage,
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

	// Parse multipart form (for file uploads)
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		// Fallback to regular form parsing if not multipart
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	var imagePath string

	// Handle image upload if present
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Process the image using upload handler logic
		uploadHandler := NewUploadHandler("uploads/images")
		path, err := uploadHandler.SaveImage(file, header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		imagePath = path
	}

	// Create task
	task, err := h.createTask.Execute(r.Context(), title, description, userID, imagePath)
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

// ShareTask handles task sharing via web form
func (h *WebTaskHandler) ShareTask(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := r.PathValue("id")

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	shareWithUserID := r.FormValue("share_with_user_id")
	if shareWithUserID == "" {
		http.Error(w, "share_with_user_id is required", http.StatusBadRequest)
		return
	}

	// Execute share use case
	err := h.shareTask.Execute(r.Context(), taskID, userID, shareWithUserID)
	if err != nil {
		if err.Error() == "only the task owner can share the task" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success message as HTML fragment
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">Tarefa compartilhada com sucesso!</div>`))
}

// DeleteTaskImage handles deleting an image from a task
func (h *WebTaskHandler) DeleteTaskImage(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := r.PathValue("id")

	// Execute delete image use case
	oldImagePath, err := h.deleteTaskImage.Execute(r.Context(), taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete the physical file
	if oldImagePath != "" {
		uploadHandler := NewUploadHandler("uploads/images")
		uploadHandler.DeleteImage(oldImagePath)
	}

	// Return empty response for HTMX to remove the image
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

// ReplaceTaskImage handles replacing an image in a task
func (h *WebTaskHandler) ReplaceTaskImage(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := r.PathValue("id")

	// Parse multipart form for image upload
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Handle new image upload
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the new image
	uploadHandler := NewUploadHandler("uploads/images")
	newImagePath, err := uploadHandler.SaveImage(file, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute replace image use case
	oldImagePath, err := h.replaceTaskImage.Execute(r.Context(), taskID, userID, newImagePath)
	if err != nil {
		// If use case fails, delete the newly uploaded image
		uploadHandler.DeleteImage(newImagePath)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete the old physical file
	if oldImagePath != "" {
		uploadHandler.DeleteImage(oldImagePath)
	}

	// Return HTML fragment with new image
	w.Header().Set("Content-Type", "text/html")
	html := `<div class="mt-3">
		<img src="` + newImagePath + `" alt="Task image" class="max-w-[200px] max-h-[200px] object-cover rounded-lg shadow-sm">
	</div>`
	w.Write([]byte(html))
}
