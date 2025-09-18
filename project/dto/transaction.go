package dto

import "time"

// TransactionDTO represents the data transfer object for Transaction
type TransactionDTO struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	AccountID string    `json:"account_id" bson:"account_id" validate:"required"`
	Amount    float64   `json:"amount" bson:"amount" validate:"required"`
	Type      string    `json:"type" bson:"type" validate:"required,oneof=deposit withdrawal"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	AccountID string  `json:"account_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required"`
	Type      string  `json:"type" validate:"required,oneof=deposit withdrawal"`
}

// TransactionResponse represents the response for transaction operations
type TransactionResponse struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransactionSummary represents a summary of transactions for an account
type TransactionSummary struct {
	AccountID         string     `json:"account_id"`
	TotalTransactions int        `json:"total_transactions"`
	TotalDeposits     float64    `json:"total_deposits"`
	TotalWithdrawals  float64    `json:"total_withdrawals"`
	CurrentBalance    float64    `json:"current_balance"`
	LastTransactionAt *time.Time `json:"last_transaction_at,omitempty"`
}
