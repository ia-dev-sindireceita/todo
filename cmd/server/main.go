package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
	"github.com/ia-edev-sindireceita/todo/internal/infrastructure/database"
	"github.com/ia-edev-sindireceita/todo/internal/infrastructure/http/handler"
	"github.com/ia-edev-sindireceita/todo/internal/infrastructure/http/middleware"
	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

func main() {
	// Initialize database
	db, err := database.NewSQLiteDB("todo.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	taskRepo := database.NewSQLiteTaskRepository(db)
	_ = database.NewSQLiteUserRepository(db) // userRepo for future use
	shareRepo := database.NewSQLiteShareRepository(db)

	// Initialize services
	taskService := service.NewTaskService(taskRepo, shareRepo)

	// Initialize use cases
	createTask := usecases.NewCreateTaskUseCase(taskRepo)
	updateTask := usecases.NewUpdateTaskUseCase(taskRepo, taskService)
	deleteTask := usecases.NewDeleteTaskUseCase(taskRepo, taskService)
	getTask := usecases.NewGetTaskUseCase(taskRepo, taskService)
	listTasks := usecases.NewListTasksUseCase(taskRepo)
	listSharedTasks := usecases.NewListSharedTasksUseCase(taskRepo)
	_ = usecases.NewShareTaskUseCase(taskRepo, shareRepo, taskService)     // shareTask for future use
	_ = usecases.NewUnshareTaskUseCase(shareRepo, taskService)            // unshareTask for future use

	// Initialize handlers
	taskHandler := handler.NewTaskHandler(
		createTask,
		updateTask,
		deleteTask,
		getTask,
		listTasks,
		listSharedTasks,
	)

	// Web handlers (for HTMX forms)
	webTaskHandler := handler.NewWebTaskHandler(createTask, deleteTask)

	// Setup router
	mux := http.NewServeMux()

	// API routes (protected)
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	apiMux.HandleFunc("GET /tasks", taskHandler.ListTasks)
	apiMux.HandleFunc("GET /tasks/shared", taskHandler.ListSharedTasks)
	apiMux.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)
	apiMux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	apiMux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	// Apply auth middleware to API routes
	mux.Handle("/api/", http.StripPrefix("/api", middleware.Chain(
		apiMux,
		middleware.AuthMiddleware,
		middleware.ContentTypeJSON,
	)))

	// Web routes (HTML)
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", handleIndex)
	webMux.HandleFunc("/tasks", handleTasksPage(listTasks))
	mux.Handle("/", webMux)

	// Web API routes (for HTMX - no auth middleware needed, uses headers)
	mux.HandleFunc("POST /web/tasks", webTaskHandler.CreateTask)
	mux.HandleFunc("DELETE /web/tasks/{id}", webTaskHandler.DeleteTask)

	// Apply global middlewares
	handler := middleware.Chain(
		mux,
		middleware.RecoverMiddleware,
		middleware.LoggingMiddleware,
		middleware.SecurityHeadersMiddleware,
		middleware.CORSMiddleware,
	)

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Database: todo.db")
	log.Println("")
	log.Println("To test the API, use:")
	log.Println("  curl -H 'X-User-ID: user-1' -H 'Content-Type: application/json' \\")
	log.Println("    -d '{\"title\":\"Test Task\",\"description\":\"Description\"}' \\")
	log.Println("    http://localhost:8080/api/tasks")
	log.Println("")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func handleTasksPage(listTasks *usecases.ListTasksUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For demo, use a hardcoded user ID
		// In production, get this from session/JWT
		userID := "user-1"

		tasks, err := listTasks.Execute(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles(
			"internal/infrastructure/templates/base.html",
			"internal/infrastructure/templates/tasks.html",
		))

		data := map[string]interface{}{
			"Title": "Tarefas",
			"Tasks": tasks,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
