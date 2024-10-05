package dto

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EntryRequest struct {
	Amount int64 `json:"amount"`
}

type TransferRequest struct {
	ToAccountID string `json:"to_account_id"`
	Amount      int64  `json:"amount"`
}

type AccountRequest struct {
	AccountID string `json:"account_id"`
	Balance   int64  `json:"balance"`
}

type LogoutRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type StoreRequest struct {
	UserID      string `json:"user_id"`
	StoreName   string `json:"store_name"`
	Description string `json:"description"`
	Email       string `json:"email"`
	NoHp        string `json:"no_hp"`
	Address     string `json:"address"`
}

type ProductRequest struct {
	StoreID     string `json:"store_id"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	Stock       int64  `json:"stock"`
	Price       int64  `json:"price"`
}

type CartRequest struct {
	CustomerID string `json:"customer_id"`
	ProductID  string `json:"product_id"`
	Quantity   int64  `json:"quantity"`
	Price      int64  `json:"price"`
	Total      int64  `json:"total"`
}