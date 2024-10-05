package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/google/uuid"
)

type CartRepository interface {
	CreateCart(newCart *dto.CartRequest) (*dto.CartResponse, error)
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
	return &newResponse, nil
}

func (repo *cartRepository) UpdateStatusCart(newStatus string, idCart string, tx *sql.Tx) (*dto.CartResponse, error) {
	stmt, err := repo.db.Prepare("UPDATE carts SET status = $1 WHERE id = $2")
	if err != nil {
		return nil, fmt.Errorf("failed to update cart : %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newStatus, idCart)

	validate(err, "update store", tx)

	return updateStore, nil
}
