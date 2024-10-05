package repository

import (
	"database/sql"
	"fmt"

	"github.com/adycahyoputro/merchant/model/entity"
)

type EntryRepository interface {
	CreateEntry(newEntry *entity.Entry, tx *sql.Tx) (*entity.Entry, error)
}

type entryRepository struct {
	db *sql.DB
}

func NewEntryRepository(db *sql.DB) EntryRepository {
	return &entryRepository{db: db}
}

func (repo *entryRepository) CreateEntry(newEntry *entity.Entry, tx *sql.Tx) (*entity.Entry, error) {
	stmt, err := repo.db.Prepare("INSERT INTO entries (id, account_id, amount, created_at) VALUES ($1,$2,$3,$4) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("failed to create entry : %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(newEntry.ID, newEntry.AccountID, newEntry.Amount, newEntry.CreatedAt).Scan(&newEntry.ID)
	validate(err, "create entry", tx)

	return newEntry, nil
}