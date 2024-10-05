package entity

import "time"

type Entry struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}