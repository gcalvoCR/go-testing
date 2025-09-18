package models

import "time"

type Account struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Balance   float64   `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
