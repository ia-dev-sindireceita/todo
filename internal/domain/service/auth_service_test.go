package service

import (
	"testing"
	"time"
)

func TestAuthService_GenerateToken(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		email     string
		secret    string
		wantError bool
	}{
		{
			name:      "should generate valid token",
			userID:    "user-123",
			email:     "user@example.com",
			secret:    "test-secret-key",
			wantError: false,
		},
		{
			name:      "should fail with empty secret",
			userID:    "user-123",
			email:     "user@example.com",
			secret:    "",
			wantError: true,
		},
		{
			name:      "should fail with empty userID",
			userID:    "",
			email:     "user@example.com",
			secret:    "test-secret-key",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := NewAuthService(tt.secret)
			token, err := authService.GenerateToken(tt.userID, tt.email, 24*time.Hour)

			if tt.wantError {
				if err == nil {
					t.Errorf("GenerateToken() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GenerateToken() unexpected error: %v", err)
				}
				if token == "" {
					t.Errorf("GenerateToken() returned empty token")
				}
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	secret := "test-secret-key"
	authService := NewAuthService(secret)

	tests := []struct {
		name      string
		setupToken func() string
		wantError bool
		wantUserID string
	}{
		{
			name: "should validate valid token",
			setupToken: func() string {
				token, _ := authService.GenerateToken("user-123", "user@example.com", 24*time.Hour)
				return token
			},
			wantError: false,
			wantUserID: "user-123",
		},
		{
			name: "should reject expired token",
			setupToken: func() string {
				token, _ := authService.GenerateToken("user-123", "user@example.com", -1*time.Hour)
				return token
			},
			wantError: true,
		},
		{
			name: "should reject invalid token",
			setupToken: func() string {
				return "invalid.token.here"
			},
			wantError: true,
		},
		{
			name: "should reject empty token",
			setupToken: func() string {
				return ""
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			claims, err := authService.ValidateToken(token)

			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateToken() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateToken() unexpected error: %v", err)
				}
				if claims.UserID != tt.wantUserID {
					t.Errorf("ValidateToken() userID = %v, want %v", claims.UserID, tt.wantUserID)
				}
			}
		})
	}
}

func TestAuthService_HashPassword(t *testing.T) {
	authService := NewAuthService("test-secret")

	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "should hash valid password",
			password:  "password123",
			wantError: false,
		},
		{
			name:      "should fail with empty password",
			password:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := authService.HashPassword(tt.password)

			if tt.wantError {
				if err == nil {
					t.Errorf("HashPassword() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("HashPassword() unexpected error: %v", err)
				}
				if hash == "" {
					t.Errorf("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Errorf("HashPassword() returned plaintext password")
				}
			}
		})
	}
}

func TestAuthService_VerifyPassword(t *testing.T) {
	authService := NewAuthService("test-secret")
	password := "password123"
	hash, _ := authService.HashPassword(password)

	tests := []struct {
		name     string
		hash     string
		password string
		wantErr  bool
	}{
		{
			name:     "should verify correct password",
			hash:     hash,
			password: password,
			wantErr:  false,
		},
		{
			name:     "should reject incorrect password",
			hash:     hash,
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "should reject empty password",
			hash:     hash,
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.VerifyPassword(tt.hash, tt.password)

			if tt.wantErr {
				if err == nil {
					t.Errorf("VerifyPassword() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("VerifyPassword() unexpected error: %v", err)
				}
			}
		})
	}
}
