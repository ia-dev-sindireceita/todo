package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
	"github.com/ia-edev-sindireceita/todo/internal/infrastructure/database"
	"github.com/ia-edev-sindireceita/todo/internal/infrastructure/http/handler"
	"github.com/ia-edev-sindireceita/todo/internal/infrastructure/http/middleware"
	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

func main() {
	// JWT secret - MUST be set via environment variable in production
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Use default only for development - NEVER in production
		jwtSecret = "development-secret-key-change-in-production"
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production!")
	}

	// Initialize database
	db, err := database.NewSQLiteDB("todo.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	taskRepo := database.NewSQLiteTaskRepository(db)
	userRepo := database.NewSQLiteUserRepository(db)
	shareRepo := database.NewSQLiteShareRepository(db)

	// Initialize services
	taskService := service.NewTaskService(taskRepo, shareRepo)

	// Initialize use cases
	createTask := usecases.NewCreateTaskUseCase(taskRepo)
	updateTask := usecases.NewUpdateTaskUseCase(taskRepo, taskService)
	deleteTask := usecases.NewDeleteTaskUseCase(taskRepo, taskService)
	completeTask := usecases.NewCompleteTaskUseCase(taskRepo, taskService)
	getTask := usecases.NewGetTaskUseCase(taskRepo, taskService)
	listTasks := usecases.NewListTasksUseCase(taskRepo)
	listSharedTasks := usecases.NewListSharedTasksUseCase(taskRepo)
	_ = usecases.NewShareTaskUseCase(taskRepo, shareRepo, taskService)     // shareTask for future use
	_ = usecases.NewUnshareTaskUseCase(shareRepo, taskService)            // unshareTask for future use

	// Auth use cases
	loginUseCase := usecases.NewLoginUseCase(userRepo, jwtSecret)
	registerUseCase := usecases.NewRegisterUseCase(userRepo, jwtSecret)

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
	webTaskHandler := handler.NewWebTaskHandler(createTask, deleteTask, completeTask)

	// Auth handlers
	authHandler := handler.NewAuthHandler(loginUseCase, registerUseCase)

	// Setup router
	mux := http.NewServeMux()

	// API routes (protected with JWT)
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
		middleware.AuthMiddleware(jwtSecret),
		middleware.ContentTypeJSON,
	)))

	// Auth API routes (no auth required)
	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /login", authHandler.Login)
	authMux.HandleFunc("POST /register", authHandler.Register)
	mux.Handle("/api/auth/", http.StripPrefix("/api/auth", middleware.Chain(
		authMux,
		middleware.ContentTypeJSON,
	)))

	// Web routes (HTML - no auth required)
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", handleIndex)
	webMux.HandleFunc("/login", handleLoginPage)
	webMux.HandleFunc("/register", handleRegisterPage)
	mux.Handle("/", webMux)

	// Web auth routes (no auth required)
	mux.HandleFunc("POST /web/auth/login", authHandler.WebLogin)
	mux.HandleFunc("POST /web/auth/register", authHandler.WebRegister)
	mux.HandleFunc("POST /web/auth/logout", authHandler.Logout)

	// Protected web routes (require JWT)
	protectedWebMux := http.NewServeMux()
	protectedWebMux.HandleFunc("/tasks", handleTasksPage(listTasks))
	mux.Handle("/tasks", middleware.AuthMiddleware(jwtSecret)(protectedWebMux))

	// Web API routes (for HTMX - require JWT)
	protectedWebAPIMux := http.NewServeMux()
	protectedWebAPIMux.HandleFunc("POST /web/tasks", webTaskHandler.CreateTask)
	protectedWebAPIMux.HandleFunc("POST /web/tasks/{id}/complete", webTaskHandler.CompleteTask)
	protectedWebAPIMux.HandleFunc("DELETE /web/tasks/{id}", webTaskHandler.DeleteTask)

	mux.Handle("/web/tasks", middleware.AuthMiddleware(jwtSecret)(protectedWebAPIMux))
	mux.Handle("/web/tasks/", middleware.AuthMiddleware(jwtSecret)(protectedWebAPIMux))

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
	http.Redirect(w, r, "/login", http.StatusFound)
}

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"internal/infrastructure/templates/base.html",
		"internal/infrastructure/templates/login.html",
	))

	data := map[string]interface{}{
		"Title": "Login",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"internal/infrastructure/templates/base.html",
		"internal/infrastructure/templates/register.html",
	))

	data := map[string]interface{}{
		"Title": "Cadastro",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleTasksPage(listTasks *usecases.ListTasksUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("userID").(string)
		if !ok || userID == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

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
			"Title":  "Tarefas",
			"Tasks":  tasks,
			"UserID": userID,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
