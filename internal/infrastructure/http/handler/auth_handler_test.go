package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// Mock for LoginUseCase
type mockLoginUseCase struct {
	executeFunc func(ctx context.Context, email, password string) (string, error)
}

func (m *mockLoginUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, email, password)
	}
	return "mock-jwt-token", nil
}

// Mock for RegisterUseCase
type mockRegisterUseCase struct {
	executeFunc func(ctx context.Context, name, email, password string) (*application.User, error)
}

func (m *mockRegisterUseCase) Execute(ctx context.Context, name, email, password string) (*application.User, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, name, email, password)
	}
	return &application.User{
		ID:        "test-user-id",
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil
}


// =============================================================================
// Login API Tests
// =============================================================================

func TestLogin_Success(t *testing.T) {
	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			if email == "test@example.com" && password == "password123" {
				return "valid-jwt-token", nil
			}
			return "", errors.New("invalid credentials")
		},
	}

	handler := &AuthHandler{loginUseCase: mockLogin}

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response LoginResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Token != "valid-jwt-token" {
		t.Errorf("Expected token 'valid-jwt-token', got %s", response.Token)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	handler := &AuthHandler{loginUseCase: &mockLoginUseCase{}}

	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Invalid request body") {
		t.Errorf("Expected error message about invalid request body, got: %s", body)
	}
}

func TestLogin_EmptyEmail(t *testing.T) {
	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			if email == "" {
				return "", errors.New("email cannot be empty")
			}
			return "token", nil
		},
	}

	handler := &AuthHandler{loginUseCase: mockLogin}

	reqBody := LoginRequest{
		Email:    "",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestLogin_EmptyPassword(t *testing.T) {
	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			if password == "" {
				return "", errors.New("password cannot be empty")
			}
			return "token", nil
		},
	}

	handler := &AuthHandler{loginUseCase: mockLogin}

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			return "", errors.New("invalid credentials")
		},
	}

	handler := &AuthHandler{loginUseCase: mockLogin}

	reqBody := LoginRequest{
		Email:    "wrong@example.com",
		Password: "wrongpass",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// =============================================================================
// Register API Tests
// =============================================================================

func TestRegister_Success(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			return &application.User{
				ID:        "new-user-id",
				Name:      name,
				Email:     email,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := &AuthHandler{registerUseCase: mockRegister}

	reqBody := RegisterRequest{
		Name:     "Test User",
		Email:    "new@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response RegisterResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != "new-user-id" {
		t.Errorf("Expected ID 'new-user-id', got %s", response.ID)
	}

	if response.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", response.Name)
	}

	if response.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got %s", response.Email)
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	handler := &AuthHandler{registerUseCase: &mockRegisterUseCase{}}

	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader("invalid-json"))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			if len(password) < 8 {
				return nil, errors.New("password must be at least 8 characters")
			}
			return &application.User{}, nil
		},
	}

	handler := &AuthHandler{registerUseCase: mockRegister}

	reqBody := RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "short",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "password must be at least 8 characters") {
		t.Errorf("Expected password validation error, got: %s", responseBody)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			if email == "existing@example.com" {
				return nil, errors.New("email already registered")
			}
			return &application.User{}, nil
		},
	}

	handler := &AuthHandler{registerUseCase: mockRegister}

	reqBody := RegisterRequest{
		Name:     "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "email already registered") {
		t.Errorf("Expected duplicate email error, got: %s", responseBody)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			return nil, errors.New("invalid email format")
		},
	}

	handler := &AuthHandler{registerUseCase: mockRegister}

	reqBody := RegisterRequest{
		Name:     "Test User",
		Email:    "invalid-email",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRegister_EmptyName(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			return nil, errors.New("user name cannot be empty")
		},
	}

	handler := &AuthHandler{registerUseCase: mockRegister}

	reqBody := RegisterRequest{
		Name:     "",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRegister_LongName(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			if len(name) > 100 {
				return nil, errors.New("user name cannot exceed 100 characters")
			}
			return &application.User{}, nil
		},
	}

	handler := &AuthHandler{registerUseCase: mockRegister}

	longName := strings.Repeat("a", 101)
	reqBody := RegisterRequest{
		Name:     longName,
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// =============================================================================
// WebLogin Tests (HTMX)
// =============================================================================

func TestWebLogin_Success(t *testing.T) {
	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			if email == "test@example.com" && password == "password123" {
				return "valid-jwt-token", nil
			}
			return "", errors.New("invalid credentials")
		},
	}

	handler := &AuthHandler{loginUseCase: mockLogin}

	formData := url.Values{}
	formData.Set("email", "test@example.com")
	formData.Set("password", "password123")

	req := httptest.NewRequest("POST", "/web/auth/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.WebLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check for HX-Redirect header
	if redirect := w.Header().Get("HX-Redirect"); redirect != "/tasks" {
		t.Errorf("Expected HX-Redirect to /tasks, got %s", redirect)
	}

	// Check for HttpOnly cookie
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == AuthCookieName {
			found = true
			if !cookie.HttpOnly {
				t.Error("Expected cookie to be HttpOnly")
			}
			if cookie.Value != "valid-jwt-token" {
				t.Errorf("Expected cookie value 'valid-jwt-token', got %s", cookie.Value)
			}
			if cookie.Path != "/" {
				t.Errorf("Expected cookie path '/', got %s", cookie.Path)
			}
		}
	}
	if !found {
		t.Error("Expected auth cookie to be set")
	}
}

func TestWebLogin_InvalidCredentials(t *testing.T) {
	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			return "", errors.New("invalid credentials")
		},
	}

	handler := &AuthHandler{loginUseCase: mockLogin}

	formData := url.Values{}
	formData.Set("email", "wrong@example.com")
	formData.Set("password", "wrongpass")

	req := httptest.NewRequest("POST", "/web/auth/login", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.WebLogin(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	// Check for error HTML fragment
	body := w.Body.String()
	if !strings.Contains(body, "Credenciais inválidas") {
		t.Errorf("Expected error HTML fragment, got: %s", body)
	}

	// Should contain Tailwind error classes
	if !strings.Contains(body, "bg-red-100") || !strings.Contains(body, "border-red-400") {
		t.Errorf("Expected Tailwind error classes in response")
	}
}

func TestWebLogin_InvalidForm(t *testing.T) {
	handler := &AuthHandler{loginUseCase: &mockLoginUseCase{}}

	// Create request with invalid form encoding
	req := httptest.NewRequest("POST", "/web/auth/login", strings.NewReader("%invalid%form%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.WebLogin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// =============================================================================
// WebRegister Tests (HTMX)
// =============================================================================

func TestWebRegister_Success(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			return &application.User{
				ID:        "new-user-id",
				Name:      name,
				Email:     email,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			return "auto-login-token", nil
		},
	}

	handler := &AuthHandler{
		registerUseCase: mockRegister,
		loginUseCase:    mockLogin,
	}

	formData := url.Values{}
	formData.Set("name", "New User")
	formData.Set("email", "new@example.com")
	formData.Set("password", "password123")

	req := httptest.NewRequest("POST", "/web/auth/register", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.WebRegister(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check for HX-Redirect header
	if redirect := w.Header().Get("HX-Redirect"); redirect != "/tasks" {
		t.Errorf("Expected HX-Redirect to /tasks, got %s", redirect)
	}

	// Check for HttpOnly cookie (auto-login)
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == AuthCookieName {
			found = true
			if cookie.Value != "auto-login-token" {
				t.Errorf("Expected auto-login cookie, got %s", cookie.Value)
			}
		}
	}
	if !found {
		t.Error("Expected auto-login cookie to be set")
	}
}

func TestWebRegister_ValidationError(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			return nil, errors.New("password must be at least 8 characters")
		},
	}

	handler := &AuthHandler{registerUseCase: mockRegister}

	formData := url.Values{}
	formData.Set("name", "Test User")
	formData.Set("email", "test@example.com")
	formData.Set("password", "short")

	req := httptest.NewRequest("POST", "/web/auth/register", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.WebRegister(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// Check for error HTML fragment with actual error message
	body := w.Body.String()
	if !strings.Contains(body, "password must be at least 8 characters") {
		t.Errorf("Expected error message in HTML fragment, got: %s", body)
	}

	// Should contain Tailwind error classes
	if !strings.Contains(body, "bg-red-100") {
		t.Errorf("Expected Tailwind error classes in response")
	}
}

func TestWebRegister_AutoLoginFails(t *testing.T) {
	mockRegister := &mockRegisterUseCase{
		executeFunc: func(ctx context.Context, name, email, password string) (*application.User, error) {
			return &application.User{
				ID:    "new-user-id",
				Name:  name,
				Email: email,
			}, nil
		},
	}

	mockLogin := &mockLoginUseCase{
		executeFunc: func(ctx context.Context, email, password string) (string, error) {
			return "", errors.New("auto-login failed")
		},
	}

	handler := &AuthHandler{
		registerUseCase: mockRegister,
		loginUseCase:    mockLogin,
	}

	formData := url.Values{}
	formData.Set("name", "Test User")
	formData.Set("email", "test@example.com")
	formData.Set("password", "password123")

	req := httptest.NewRequest("POST", "/web/auth/register", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.WebRegister(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Should redirect to login page when auto-login fails
	if redirect := w.Header().Get("HX-Redirect"); redirect != "/login" {
		t.Errorf("Expected HX-Redirect to /login when auto-login fails, got %s", redirect)
	}
}

func TestWebRegister_InvalidForm(t *testing.T) {
	handler := &AuthHandler{registerUseCase: &mockRegisterUseCase{}}

	req := httptest.NewRequest("POST", "/web/auth/register", strings.NewReader("%invalid%form%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.WebRegister(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// =============================================================================
// Logout Tests
// =============================================================================

func TestLogout_Success(t *testing.T) {
	handler := &AuthHandler{}

	req := httptest.NewRequest("POST", "/web/auth/logout", nil)
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check for HX-Redirect header
	if redirect := w.Header().Get("HX-Redirect"); redirect != "/login" {
		t.Errorf("Expected HX-Redirect to /login, got %s", redirect)
	}

	// Check that cookie is deleted (MaxAge should be negative)
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == AuthCookieName {
			found = true
			if cookie.Value != "" {
				t.Errorf("Expected empty cookie value, got %s", cookie.Value)
			}
			if cookie.MaxAge >= 0 {
				t.Errorf("Expected negative MaxAge to delete cookie, got %d", cookie.MaxAge)
			}
		}
	}
	if !found {
		t.Error("Expected auth cookie to be set for deletion")
	}
}
