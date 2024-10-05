package entity

import "time"

type Store struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	StoreName 	string    `json:"store_name"`
	Description string    `json:"description"`
	Email		string    `json:"email"`
	NoHp        string    `json:"no_hp"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
	IsDelete    bool      `json:"is_delete"`
}