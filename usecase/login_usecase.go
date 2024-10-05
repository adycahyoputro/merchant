package usecase

import (
	"errors"
	"time"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/repository"
	"github.com/adycahyoputro/merchant/utils"
	"github.com/dgrijalva/jwt-go"
)

type LoginUsecase interface {
	Login(newLogin *dto.LoginRequest) (string, error)
	Logout(userID string) error
}

type loginUsecase struct {
	userRepo    repository.UserRepository
	accountRepo repository.AccountRepository
}

func NewLoginUsecase(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository) LoginUsecase {
	return &loginUsecase{
		userRepo:    userRepo,
		accountRepo: accountRepo}
}

func (usecase *loginUsecase) Login(newLogin *dto.LoginRequest) (string, error) {
	if newLogin.Email == "" {
		return "", errors.New("email is required")
	}
	if newLogin.Password == "" {
		return "", errors.New("password is required")
	}
	getUserByEmail, _ := usecase.userRepo.FindUserByEmail(newLogin.Email)
	if getUserByEmail == nil {
		return "", errors.New("user with email " + newLogin.Email + " not registered")
	}

	match := utils.CheckPasswordHash(newLogin.Password, getUserByEmail.Password)
	if !match {
		return "", errors.New("wrong password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = getUserByEmail.ID
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	var jwtKeyByte = []byte(utils.GetEnv("JWT_KEY"))
	tokenString, err := token.SignedString(jwtKeyByte)
	if err != nil {
		return "", errors.New("internal server error")
	}

	newUpdateAccount := dto.LogoutRequest{
		UserID:   getUserByEmail.ID,
		IsActive: true,
	}
	err = usecase.accountRepo.UpdateIsActiveAccount(&newUpdateAccount)
	if err != nil {
		return "", errors.New("internal server error")
	}

	return tokenString, nil
}

func (usecase *loginUsecase) Logout(userID string) error {

	newUpdateAccount := dto.LogoutRequest{
		UserID:   userID,
		IsActive: false,
	}
	err := usecase.accountRepo.UpdateIsActiveAccount(&newUpdateAccount)
	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}