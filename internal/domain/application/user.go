package application

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// User represents a user entity
type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewUser creates a new User with validation
func NewUser(id, name, email, passwordHash string) (*User, error) {
	if id == "" {
		return nil, errors.New("user id cannot be empty")
	}

	if name == "" {
		return nil, errors.New("user name cannot be empty")
	}

	if len(name) > 100 {
		return nil, errors.New("user name cannot exceed 100 characters")
	}

	if email == "" {
		return nil, errors.New("user email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	if passwordHash == "" {
		return nil, errors.New("password hash cannot be empty")
	}

	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}, nil
}
