package repository

import (
	"database/sql"
	"fmt"

	"github.com/adycahyoputro/merchant/model/entity"
)

type TransferRepository interface {
	CreateTransfer(newTransfer *entity.Transfer, tx *sql.Tx) (*entity.Transfer, error)
}

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(
	db *sql.DB) TransferRepository {
	return &transferRepository{
		db: db}
}

func (repo *transferRepository) CreateTransfer(newTransfer *entity.Transfer, tx *sql.Tx) (*entity.Transfer, error) {
	stmt, err := repo.db.Prepare("INSERT INTO transfers (id, from_account_id, to_account_id, amount, created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer : %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(newTransfer.ID, newTransfer.FromAccountID, newTransfer.ToAccountID, newTransfer.Amount, newTransfer.CreatedAt).Scan(&newTransfer.ID)

	validate(err, "create transfer", tx)

	return newTransfer, nil
}