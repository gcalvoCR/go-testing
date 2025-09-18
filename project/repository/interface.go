package repository

import (
	"context"

	"github.com/gcalvocr/go-testing/dto"
)

// AccountRepository defines the interface for account data operations
type AccountRepository interface {
	Create(ctx context.Context, account *dto.AccountDTO) error
	GetByID(ctx context.Context, id string) (*dto.AccountDTO, error)
	GetAll(ctx context.Context) ([]*dto.AccountDTO, error)
	Update(ctx context.Context, id string, account *dto.UpdateAccountRequest) error
	Delete(ctx context.Context, id string) error
	GetByName(ctx context.Context, name string) (*dto.AccountDTO, error)
	UpdateBalance(ctx context.Context, id string, newBalance float64) error
}

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	Create(ctx context.Context, transaction *dto.TransactionDTO) error
	GetByID(ctx context.Context, id string) (*dto.TransactionDTO, error)
	GetByAccountID(ctx context.Context, accountID string) ([]*dto.TransactionDTO, error)
	GetAll(ctx context.Context) ([]*dto.TransactionDTO, error)
	Update(ctx context.Context, id string, transaction *dto.TransactionDTO) error
	Delete(ctx context.Context, id string) error
	GetTransactionSummary(ctx context.Context, accountID string) (*dto.TransactionSummary, error)
}

// DatabaseType represents the type of database
type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgres"
	MongoDB    DatabaseType = "mongodb"
)

// RepositoryFactory creates repository instances based on database type
type RepositoryFactory struct {
	AccountRepo     AccountRepository
	TransactionRepo TransactionRepository
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(dbType DatabaseType, connectionString string) (*RepositoryFactory, error) {
	switch dbType {
	case PostgreSQL:
		return newPostgreSQLFactory(connectionString)
	case MongoDB:
		return newMongoDBFactory(connectionString)
	default:
		return newPostgreSQLFactory(connectionString) // default to PostgreSQL
	}
}
