package entity

import "time"

type Carts struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	ProductID  string    `json:"product_id"`
	Quantity   int64     `json:"quantity"`
	Price      int64     `json:"price"`
	Total      int64     `json:"total"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdateAt   time.Time `json:"updated_at"`
	IsDelete   bool      `json:"is_delete"`
}