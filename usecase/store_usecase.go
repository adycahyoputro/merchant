package usecase

import (
	"errors"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/repository"
)

type StoreUsecase interface {
	CreateStore(newStore *dto.StoreRequest) (*dto.StoreResponse, error)
}

type storeUsecase struct {
	storeRepo repository.StoreRepository
}

func NewStoreUsecase(storeRepo repository.StoreRepository) StoreUsecase {
	return &storeUsecase{storeRepo: storeRepo}
}

func (usecase *storeUsecase) CreateStore(newStore *dto.StoreRequest) (*dto.StoreResponse, error) {
	if newStore.StoreName == "" {
		return nil, errors.New("store name is required")
	}
	if newStore.Description == "" {
		return nil, errors.New("description is required")
	}
	if newStore.Email == "" {
		return nil, errors.New("email is required")
	}
	if newStore.NoHp == "" {
		return nil, errors.New("no hp is required")
	}
	if newStore.Address == "" {
		return nil, errors.New("address is required")
	}
	if newStore.StoreName == "" {
		return nil, errors.New("store name is required")
	}
	if len(newStore.NoHp) < 11 || len(newStore.NoHp) > 12 {
		return nil, errors.New("length no hp " + newStore.NoHp + " at least 12")
	}
	return usecase.storeRepo.CreateStore(newStore)
}