package usecase

import (
	"errors"
	"fmt"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/repository"
	"golang.org/x/crypto/bcrypt"
)

type TransactionUsecase interface {
	CreateUserAccount(newMainUser *dto.UserRequest) (*dto.UserResponse, error)
	CreateMainTransfer(newMainTransfer *dto.TransferRequest) (*dto.TransferResponse, error)
	CreateMainEntry(newMainEntry *dto.EntryRequest, accountID string, balance int64) (*dto.EntryResponse, error)
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
	accountRepo     repository.AccountRepository
	cartRepo  repository.CartRepository
	productRepo repository.ProductRepository
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	cartRepo  repository.CartRepository,
	productRepo repository.ProductRepository) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		accountRepo:     accountRepo,
		cartRepo: cartRepo,
		productRepo: productRepo,}
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

func (usecase *transactionUsecase) CreateMainTransfer(newMainTransfer *dto.TransferRequest) (*dto.TransferResponse, error) {
	if newMainTransfer.ToAccountID == "" {
		return nil, errors.New("destination account is required")
	}
	// if newMainTransfer.Amount == 0 {
	// 	return nil, errors.New("amount is required")
	// }
	if newMainTransfer.Amount < 0 {
		return nil, errors.New("amount must be positive amount")
	}

	getAccountByToAccountID, err := usecase.accountRepo.FindAccountByAccountID(newMainTransfer.ToAccountID)
	if err != nil {
		return nil, errors.New("account with account id " + newMainTransfer.ToAccountID + " not found")
	}
	newMainTransfer.NewToBalance = getAccountByToAccountID.Balance + newMainTransfer.Amount

	newMainTransfer.NewStatusCart = "payment"

	getCartByUserID, err := usecase.cartRepo.FindCartByUserID(newMainTransfer.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart list : %w", err)
	}

	var newTotal int64
	var stocks = make([]int64, len(getCartByUserID))
	for _, v := range getCartByUserID {
		getProductByProductID, err := usecase.productRepo.FindProductByProductID(v.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product : %w", err)
		}
		newStock := getProductByProductID.Stock - v.Quantity
		if newStock < 0 {
			return nil, errors.New("stock is not enough")
		}
		stocks = append(stocks,newStock)
		newTotal += v.Total
		fmt.Println("total:",newTotal)
	}

	newMainTransfer.NewStock = stocks

	getAccountByFromAccountID, err := usecase.accountRepo.FindAccountByAccountID(newMainTransfer.FromAccountID)
	if err != nil {
		return nil, errors.New("account with account id " + newMainTransfer.FromAccountID + " not found")
	}
	if !getAccountByFromAccountID.IsActive {
		return nil, errors.New("user unauthorize")
	}
	// fmt.Println(newTotal)
	newMainTransfer.NewFromBalance = getAccountByFromAccountID.Balance - newTotal
	if newMainTransfer.NewFromBalance < 0 {
		return nil, errors.New("balance is not enough")
	}
	newMainTransfer.Amount = newTotal

	return usecase.transactionRepo.CreateMainTransfer(newMainTransfer)
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