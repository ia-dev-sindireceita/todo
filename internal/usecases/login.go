package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

// LoginUseCase handles user login
type LoginUseCase struct {
	userRepo    repository.UserRepository
	authService *service.AuthService
}

// NewLoginUseCase creates a new LoginUseCase
func NewLoginUseCase(userRepo repository.UserRepository, jwtSecret string) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		authService: service.NewAuthService(jwtSecret),
	}
}

// Execute performs user login and returns a JWT token
func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", errors.New("email cannot be empty")
	}
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verify password
	if err := uc.authService.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := uc.authService.GenerateToken(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}
