package application

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		userName     string
		email        string
		passwordHash string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid user",
			id:           "user-1",
			userName:     "John Doe",
			email:        "john@example.com",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMye",
			wantErr:      false,
		},
		{
			name:         "empty id",
			id:           "",
			userName:     "John Doe",
			email:        "john@example.com",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMye",
			wantErr:      true,
			errMsg:       "user id cannot be empty",
		},
		{
			name:         "empty name",
			id:           "user-1",
			userName:     "",
			email:        "john@example.com",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMye",
			wantErr:      true,
			errMsg:       "user name cannot be empty",
		},
		{
			name:         "name too long",
			id:           "user-1",
			userName:     string(make([]byte, 101)),
			email:        "john@example.com",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMye",
			wantErr:      true,
			errMsg:       "user name cannot exceed 100 characters",
		},
		{
			name:         "empty email",
			id:           "user-1",
			userName:     "John Doe",
			email:        "",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMye",
			wantErr:      true,
			errMsg:       "user email cannot be empty",
		},
		{
			name:         "invalid email format",
			id:           "user-1",
			userName:     "John Doe",
			email:        "invalid-email",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMye",
			wantErr:      true,
			errMsg:       "invalid email format",
		},
		{
			name:         "empty password hash",
			id:           "user-1",
			userName:     "John Doe",
			email:        "john@example.com",
			passwordHash: "",
			wantErr:      true,
			errMsg:       "password hash cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.userName, tt.email, tt.passwordHash)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("NewUser() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}

			if user.ID != tt.id {
				t.Errorf("User.ID = %v, want %v", user.ID, tt.id)
			}
			if user.Name != tt.userName {
				t.Errorf("User.Name = %v, want %v", user.Name, tt.userName)
			}
			if user.Email != tt.email {
				t.Errorf("User.Email = %v, want %v", user.Email, tt.email)
			}
			if user.PasswordHash != tt.passwordHash {
				t.Errorf("User.PasswordHash = %v, want %v", user.PasswordHash, tt.passwordHash)
			}
			if user.CreatedAt.IsZero() {
				t.Error("User.CreatedAt should not be zero")
			}
		})
	}
}
