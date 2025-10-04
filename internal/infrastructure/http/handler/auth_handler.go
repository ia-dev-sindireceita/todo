package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	loginUseCase    usecases.LoginUseCaseInterface
	registerUseCase usecases.RegisterUseCaseInterface
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	loginUseCase usecases.LoginUseCaseInterface,
	registerUseCase usecases.RegisterUseCaseInterface,
) *AuthHandler {
	return &AuthHandler{
		loginUseCase:    loginUseCase,
		registerUseCase: registerUseCase,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse represents a registration response
type RegisterResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Login handles user login (API)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.loginUseCase.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// Register handles user registration (API)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.registerUseCase.Execute(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

// WebLogin handles web login (form submission)
func (h *AuthHandler) WebLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, err := h.loginUseCase.Execute(r.Context(), email, password)
	if err != nil {
		// Return error HTML fragment for HTMX
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
			Credenciais inválidas. Tente novamente.
		</div>`))
		return
	}

	// Set JWT token in HttpOnly cookie
	http.SetCookie(w, createAuthCookie(token))

	// Redirect to tasks page
	w.Header().Set("HX-Redirect", "/tasks")
	w.WriteHeader(http.StatusOK)
}

// WebRegister handles web registration (form submission)
func (h *AuthHandler) WebRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.registerUseCase.Execute(r.Context(), name, email, password)
	if err != nil {
		// Return error HTML fragment for HTMX
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
			` + err.Error() + `
		</div>`))
		return
	}

	// Auto-login after registration using the same password
	token, err := h.loginUseCase.Execute(r.Context(), user.Email, password)
	if err != nil {
		// Redirect to login page if auto-login fails
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set JWT token in HttpOnly cookie
	http.SetCookie(w, createAuthCookie(token))

	// Redirect to tasks page
	w.Header().Set("HX-Redirect", "/tasks")
	w.WriteHeader(http.StatusOK)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, deleteAuthCookie())

	// Redirect to login page
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}
