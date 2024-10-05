package dto

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type UserResponse struct {
	ID       string `json:"account_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Balance  int64  `json:"balance"`
}

type TransferResponse struct {
	ID            string `json:"id"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

type EntryResponse struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
	Balance   int64  `json:"balance"`
}

type StoreResponse struct {
	ID          string `json:"store_id"`
	UserID      string `json:"user_id"`
	StoreName   string `json:"store_name"`
	Description string `json:"description"`
	Email       string `json:"email"`
	NoHp        string `json:"no_hp"`
	Address     string `json:"address"`
}

type ProductResponse struct {
	ID          string `json:"product_id"`
	StoreID     string `json:"store_id"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	Stock       int64  `json:"stock"`
	Price       int64  `json:"price"`
}

type CartResponse struct {
	ID         string `json:"cart_id"`
	CustomerID string `json:"customer_id"`
	ProductID  string `json:"product_id"`
	Quantity   int64  `json:"quantity"`
	Price      int64  `json:"price"`
	Total      int64  `json:"total"`
}