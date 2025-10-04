package usecases

import (
	"context"
	"testing"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// Mock UserRepository for testing
type mockUserRepositoryForLogin struct {
	users map[string]*application.User
}

func (m *mockUserRepositoryForLogin) Create(ctx context.Context, user *application.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepositoryForLogin) FindByID(ctx context.Context, id string) (*application.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, application.ErrUserNotFound
}

func (m *mockUserRepositoryForLogin) FindByEmail(ctx context.Context, email string) (*application.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, application.ErrUserNotFound
}

func (m *mockUserRepositoryForLogin) Update(ctx context.Context, user *application.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepositoryForLogin) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func TestLoginUseCase_Execute(t *testing.T) {
	// Setup
	mockRepo := &mockUserRepositoryForLogin{
		users: make(map[string]*application.User),
	}

	loginUseCase := NewLoginUseCase(mockRepo, "test-secret-key")

	// Create test user with properly hashed password
	// We need to hash the password using the same auth service
	passwordHash, err := loginUseCase.authService.HashPassword("password123")
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	testUser := &application.User{
		ID:           "user-1",
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
	}
	mockRepo.users[testUser.ID] = testUser

	tests := []struct {
		name      string
		email     string
		password  string
		wantError bool
	}{
		{
			name:      "should login with correct credentials",
			email:     "test@example.com",
			password:  "password123",
			wantError: false,
		},
		{
			name:      "should fail with incorrect password",
			email:     "test@example.com",
			password:  "wrongpassword",
			wantError: true,
		},
		{
			name:      "should fail with non-existent email",
			email:     "nonexistent@example.com",
			password:  "password123",
			wantError: true,
		},
		{
			name:      "should fail with empty email",
			email:     "",
			password:  "password123",
			wantError: true,
		},
		{
			name:      "should fail with empty password",
			email:     "test@example.com",
			password:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := loginUseCase.Execute(context.Background(), tt.email, tt.password)

			if tt.wantError {
				if err == nil {
					t.Errorf("Execute() expected error but got nil")
				}
				if token != "" {
					t.Errorf("Execute() expected empty token on error")
				}
			} else {
				if err != nil {
					t.Errorf("Execute() unexpected error: %v", err)
				}
				if token == "" {
					t.Errorf("Execute() expected token but got empty string")
				}
			}
		})
	}
}
