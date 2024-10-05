package entity

import "time"

type Product struct {
	ID        string    `json:"id"`
	StoreID   string    `json:"store_id"`
	ProductName     string    `json:"product_name"`
	Description		string		`json:"description"`
	Stock   int64     `json:"stock"`
	Price  int64    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
	IsDelete  bool      `json:"is_delete"`
}