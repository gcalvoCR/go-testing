package main

import (
	"errors"
	"fmt"
)

// User represents a user in the system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Database interface for user operations
type UserDatabase interface {
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id int) error
}

// UserService provides user-related business logic
type UserService struct {
	db UserDatabase
}

// NewUserService creates a new user service
func NewUserService(db UserDatabase) *UserService {
	return &UserService{db: db}
}

// GetUser retrieves a user by ID with business logic
func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.db.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	user := &User{
		Name:  name,
		Email: email,
	}

	err := s.db.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUserEmail updates a user's email
func (s *UserService) UpdateUserEmail(id int, newEmail string) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}
	if newEmail == "" {
		return errors.New("email is required")
	}

	// First get the user
	user, err := s.db.GetUserByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user for update: %w", err)
	}

	// Update email
	user.Email = newEmail

	// Save changes
	err = s.db.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
