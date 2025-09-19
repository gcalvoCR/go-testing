package repository

import (
	"context"
	"testcontainers-mongo-demo/models"
)

type UserRepository interface {
	Ping(ctx context.Context) error
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*models.User, error)
	Close(ctx context.Context) error
}

// DatabaseType represents the type of database
type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgres" // not implemented
	MongoDB    DatabaseType = "mongodb"
)
