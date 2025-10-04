package usecases

import (
	"context"
	"testing"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// Mock UserRepository for testing
type mockUserRepositoryForRegister struct {
	users map[string]*application.User
}

func (m *mockUserRepositoryForRegister) Create(ctx context.Context, user *application.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepositoryForRegister) FindByID(ctx context.Context, id string) (*application.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, application.ErrUserNotFound
}

func (m *mockUserRepositoryForRegister) FindByEmail(ctx context.Context, email string) (*application.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, application.ErrUserNotFound
}

func (m *mockUserRepositoryForRegister) Update(ctx context.Context, user *application.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepositoryForRegister) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func TestRegisterUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		userName  string
		email     string
		password  string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "should register new user with valid data",
			userName:  "John Doe",
			email:     "john@example.com",
			password:  "password123",
			wantError: false,
		},
		{
			name:      "should fail with empty name",
			userName:  "",
			email:     "john@example.com",
			password:  "password123",
			wantError: true,
		},
		{
			name:      "should fail with invalid email",
			userName:  "John Doe",
			email:     "invalid-email",
			password:  "password123",
			wantError: true,
		},
		{
			name:      "should fail with empty email",
			userName:  "John Doe",
			email:     "",
			password:  "password123",
			wantError: true,
		},
		{
			name:      "should fail with empty password",
			userName:  "John Doe",
			email:     "john@example.com",
			password:  "",
			wantError: true,
		},
		{
			name:      "should fail with short password",
			userName:  "John Doe",
			email:     "john@example.com",
			password:  "short",
			wantError: true,
			errorMsg:  "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepositoryForRegister{
				users: make(map[string]*application.User),
			}
			registerUseCase := NewRegisterUseCase(mockRepo, "test-secret-key")

			user, err := registerUseCase.Execute(context.Background(), tt.userName, tt.email, tt.password)

			if tt.wantError {
				if err == nil {
					t.Errorf("Execute() expected error but got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Execute() error = %v, want %v", err.Error(), tt.errorMsg)
				}
				if user != nil {
					t.Errorf("Execute() expected nil user on error")
				}
			} else {
				if err != nil {
					t.Errorf("Execute() unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("Execute() expected user but got nil")
				}
				if user != nil {
					if user.Name != tt.userName {
						t.Errorf("Execute() user.Name = %v, want %v", user.Name, tt.userName)
					}
					if user.Email != tt.email {
						t.Errorf("Execute() user.Email = %v, want %v", user.Email, tt.email)
					}
					if user.PasswordHash == "" || user.PasswordHash == tt.password {
						t.Errorf("Execute() password not properly hashed")
					}
				}
			}
		})
	}
}

func TestRegisterUseCase_Execute_DuplicateEmail(t *testing.T) {
	mockRepo := &mockUserRepositoryForRegister{
		users: make(map[string]*application.User),
	}
	registerUseCase := NewRegisterUseCase(mockRepo, "test-secret-key")

	// Register first user
	_, err := registerUseCase.Execute(context.Background(), "User One", "duplicate@example.com", "password123")
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to register with same email
	_, err = registerUseCase.Execute(context.Background(), "User Two", "duplicate@example.com", "password456")
	if err == nil {
		t.Errorf("Execute() expected error for duplicate email but got nil")
	}
	if err != nil && err.Error() != "email already registered" {
		t.Errorf("Execute() error = %v, want 'email already registered'", err.Error())
	}
}
