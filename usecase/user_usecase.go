package usecase

import (
	"errors"

	"github.com/adycahyoputro/merchant/model/entity"
	"github.com/adycahyoputro/merchant/repository"
)

type UserUsecase interface {
	FindUserByEmail(email string) (*entity.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (usecase *userUsecase) FindUserByEmail(email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return usecase.userRepo.FindUserByEmail(email)
}