package usecase

import (
	"errors"

	"github.com/adycahyoputro/merchant/model/entity"
	"github.com/adycahyoputro/merchant/repository"
)

type AccountUsecase interface {
	FindAccountByUserID(userID string) (*entity.Account, error)
}

type accountUsecase struct {
	accountRepo repository.AccountRepository
}

func NewAccountUsecase(accountRepo repository.AccountRepository) AccountUsecase {
	return &accountUsecase{accountRepo: accountRepo}
}

func (usecase *accountUsecase) FindAccountByUserID(userID string) (*entity.Account, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	return usecase.accountRepo.FindAccountByUserID(userID)
}