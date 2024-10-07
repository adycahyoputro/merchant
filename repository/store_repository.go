package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/google/uuid"
)

type StoreRepository interface {
	CreateStore(newStore *dto.StoreRequest) (*dto.StoreResponse, error)
}

type storeRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (repo *storeRepository) CreateStore(newStore *dto.StoreRequest) (*dto.StoreResponse, error) {
	stmt, err := repo.db.Prepare("INSERT INTO stores (id, user_id, store_name, description, email, no_hp, address, created_at, updated_at, is_delete) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("failed to create store : %w", err)
	}
	defer stmt.Close()

	repeat := true
	var idStore string
	for repeat {
		idStoreUID, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to create uuid : %w", err)
		}
		if idStoreUID.String() != newStore.UserID {
			repeat = false
			idStore = idStoreUID.String()
		}
	}
	
	createdAt := time.Now()
	updatedAt := time.Now()
	isDelete := false
	err = stmt.QueryRow(idStore, newStore.UserID, newStore.StoreName, newStore.Description, newStore.Email, newStore.NoHp, newStore.Address, createdAt, updatedAt, isDelete).Scan(idStore)
	if err != nil {
		return nil, fmt.Errorf("failed to create store : %w", err)
	}
	newResponse := dto.StoreResponse{
		ID: idStore,
		UserID: newStore.UserID,
		StoreName: newStore.StoreName,
		Description: newStore.Description,
		Email: newStore.Email,
		NoHp: newStore.NoHp,
		Address: newStore.Address,
	}
	fmt.Println(newResponse)
	return &newResponse, nil
}