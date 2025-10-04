package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/ia-edev-sindireceita/todo/internal/domain/service"
)

// RegisterUseCase handles user registration
type RegisterUseCase struct {
	userRepo    repository.UserRepository
	authService *service.AuthService
}

// NewRegisterUseCase creates a new RegisterUseCase
func NewRegisterUseCase(userRepo repository.UserRepository, jwtSecret string) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:    userRepo,
		authService: service.NewAuthService(jwtSecret),
	}
}

// Execute registers a new user
func (uc *RegisterUseCase) Execute(ctx context.Context, name, email, password string) (*application.User, error) {
	// Validate password length
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if email already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	passwordHash, err := uc.authService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	id := uuid.New().String()
	user, err := application.NewUser(id, name, email, passwordHash)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
