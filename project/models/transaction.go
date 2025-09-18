package models

import "time"

type Transaction struct {
	ID        int       `json:"id" db:"id"`
	AccountID int       `json:"account_id" db:"account_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Type      string    `json:"type" db:"type"` // deposit, withdrawal
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
