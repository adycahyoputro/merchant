package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/model/entity"
	"github.com/google/uuid"
)

type CartRepository interface {
	CreateCart(newCart *dto.CartRequest) (*dto.CartResponse, error)
	UpdateStatusCart(newStatus string, idCart string, tx *sql.Tx) error
	FindCartByUserID(userID string) ([]entity.Carts, error)
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (repo *cartRepository) CreateCart(newCart *dto.CartRequest) (*dto.CartResponse, error) {
	stmt, err := repo.db.Prepare("INSERT INTO carts (id, customer_id, product_id, quantity, price, total, status, created_at, updated_at, is_delete) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("failed to create cart : %w", err)
	}
	defer stmt.Close()

	repeat := true
	var idCart string
	for repeat {
		idCartUID, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to create uuid : %w", err)
		}
		if idCartUID.String() != newCart.CustomerID {
			repeat = false
			idCart = idCartUID.String()
		}
	}

	status := "pending"
	createdAt := time.Now()
	updatedAt := time.Now()
	isDelete := false
	err = stmt.QueryRow(idCart, newCart.CustomerID, newCart.ProductID, newCart.Quantity, newCart.Price, newCart.Total, status, createdAt, updatedAt, isDelete).Scan(idCart)
	if err != nil {
		return nil, fmt.Errorf("failed to create store : %w", err)
	}
	newResponse := dto.CartResponse{
		ID:          idCart,
		CustomerID: newCart.CustomerID,
		ProductID: newCart.ProductID,
		Quantity: newCart.Quantity,
		Price: newCart.Price,
		Total: newCart.Total,
	}
	fmt.Println(newResponse)
	return &newResponse, nil
}

func (repo *cartRepository) UpdateStatusCart(newStatus string, idCart string, tx *sql.Tx) error {
	stmt, err := repo.db.Prepare("UPDATE carts SET status = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("failed to update cart : %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newStatus, idCart)

	validate(err, "update cart", tx)

	return nil
}

func (repo *cartRepository) FindCartByUserID(userID string) ([]entity.Carts, error) {
	var carts []entity.Carts
	rows, err := repo.db.Query("SELECT id, customer_id, product_id, quantity, price, total, status, created_at, updated_at, is_delete FROM carts WHERE customer_id = $1 AND status = 'pending' ORDER BY id", userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart with user id %v not found", userID)
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cart entity.Carts
		err := rows.Scan(&cart.ID, &cart.CustomerID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Total, &cart.Status, &cart.CreatedAt, &cart.UpdateAt, &cart.IsDelete)
		if err != nil {
			return nil, fmt.Errorf("failed to get cart : %w", err)
		}
		carts = append(carts, cart)
	}

	return carts, nil
}
