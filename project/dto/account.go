package dto

import "time"

// AccountDTO represents the data transfer object for Account
type AccountDTO struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name" validate:"required,min=1,max=100"`
	Balance   float64   `json:"balance" bson:"balance" validate:"min=0"`
	Currency  string    `json:"currency" bson:"currency" validate:"required,len=3"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CreateAccountRequest represents the request to create an account
type CreateAccountRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=100"`
	Balance  float64 `json:"balance" validate:"min=0"`
	Currency string  `json:"currency" validate:"required,len=3"`
}

// UpdateAccountRequest represents the request to update an account
type UpdateAccountRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Balance  *float64 `json:"balance,omitempty" validate:"omitempty,min=0"`
	Currency *string  `json:"currency,omitempty" validate:"omitempty,len=3"`
}

// AccountResponse represents the response for account operations
type AccountResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
