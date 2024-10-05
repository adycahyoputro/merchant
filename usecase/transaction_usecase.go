package usecase

import (
	"errors"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/repository"
	"golang.org/x/crypto/bcrypt"
)

type TransactionUsecase interface {
	CreateUserAccount(newMainUser *dto.UserRequest) (*dto.UserResponse, error)
	CreateMainTransfer(newMainTransfer *dto.TransferRequest, fromAccountID string) (*dto.TransferResponse, error)
	CreateMainEntry(newMainEntry *dto.EntryRequest, accountID string, balance int64) (*dto.EntryResponse, error)
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
	accountRepo     repository.AccountRepository
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		accountRepo:     accountRepo}
}

func (usecase *transactionUsecase) CreateUserAccount(newUserAccount *dto.UserRequest) (*dto.UserResponse, error) {
	if newUserAccount.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if newUserAccount.LastName == "" {
		return nil, errors.New("last name is required")
	}
	if newUserAccount.Email == "" {
		return nil, errors.New("email is required")
	}
	if newUserAccount.Password == "" {
		return nil, errors.New("password is required")
	}

	getUserByEmail, _ := usecase.userRepo.FindUserByEmail(newUserAccount.Email)
	if getUserByEmail != nil {
		return nil, errors.New("user with email " + getUserByEmail.Email + " has registered")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newUserAccount.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUserAccount.Password = string(encryptedPassword)

	return usecase.transactionRepo.CreateUserAccount(newUserAccount)
}

func (usecase *transactionUsecase) CreateMainTransfer(newMainTransfer *dto.TransferRequest, fromAccountID string) (*dto.TransferResponse, error) {
	if newMainTransfer.ToAccountID == "" {
		return nil, errors.New("destination account is required")
	}
	if newMainTransfer.Amount == 0 {
		return nil, errors.New("amount is required")
	}
	if newMainTransfer.Amount < 0 {
		return nil, errors.New("amount must be positive amount")
	}

	getAccountByFromAccountID, err := usecase.accountRepo.FindAccountByAccountID(fromAccountID)
	if err != nil {
		return nil, errors.New("account with account id " + fromAccountID + " not found")
	}
	if !getAccountByFromAccountID.IsActive {
		return nil, errors.New("user unauthorize")
	}
	newFromBalance := getAccountByFromAccountID.Balance - newMainTransfer.Amount
	if newFromBalance < 0 {
		return nil, errors.New("balance is not enough")
	}

	getAccountByToAccountID, err := usecase.accountRepo.FindAccountByAccountID(newMainTransfer.ToAccountID)
	if err != nil {
		return nil, errors.New("account with account id " + newMainTransfer.ToAccountID + " not found")
	}
	newToBalance := getAccountByToAccountID.Balance + newMainTransfer.Amount

	return usecase.transactionRepo.CreateMainTransfer(newMainTransfer, fromAccountID, newFromBalance, newToBalance)
}

func (usecase *transactionUsecase) CreateMainEntry(newMainEntry *dto.EntryRequest, accountID string, balance int64) (*dto.EntryResponse, error) {
	if newMainEntry.Amount == 0 {
		return nil, errors.New("amount is required")
	}
	if newMainEntry.Amount < 0 {
		return nil, errors.New("amount must be positive amount")
	}
	getAccountByFromAccountID, err := usecase.accountRepo.FindAccountByAccountID(accountID)
	if err != nil {
		return nil, errors.New("account with account id " + accountID + " not found")
	}
	if !getAccountByFromAccountID.IsActive {
		return nil, errors.New("user unauthorize")
	}
	newBalance := balance + newMainEntry.Amount

	return usecase.transactionRepo.CreateMainEntry(newMainEntry, accountID, newBalance)
}