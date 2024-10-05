package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/model/entity"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateUserAccount(newMainUser *dto.UserRequest) (*dto.UserResponse, error)
	CreateMainTransfer(newMainTransfer *dto.TransferRequest, fromAccountID string, fromBalance int64, toBalance int64) (*dto.TransferResponse, error)
	CreateMainEntry(newMainEntry *dto.EntryRequest, accountID string, balance int64) (*dto.EntryResponse, error)
}

type transactionRepository struct {
	db                 *sql.DB
	userRepository     UserRepository
	accountRepository  AccountRepository
	entryRepository    EntryRepository
	transferRepository TransferRepository
}

func NewTransactionRepository(
	db *sql.DB,
	userRepository UserRepository,
	accountRepository AccountRepository,
	entryRepository EntryRepository,
	transferRepository TransferRepository) TransactionRepository {
	return &transactionRepository{
		db:                 db,
		userRepository:     userRepository,
		accountRepository:  accountRepository,
		entryRepository:    entryRepository,
		transferRepository: transferRepository,
	}
}

func (repo *transactionRepository) CreateUserAccount(newMainUser *dto.UserRequest) (*dto.UserResponse, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create user : %w", err)
	}

	idUser, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid : %w", err)
	}
	repeat := true
	var idAccount string
	for repeat {
		idAccountUID, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to create uuid : %w", err)
		}
		if idAccountUID.String() != idUser.String() {
			repeat = false
			idAccount = idAccountUID.String()
		}
	}
	createdAt := time.Now()

	newUser := entity.User{
		ID:        idUser.String(),
		FirstName: newMainUser.FirstName,
		LastName:  newMainUser.LastName,
		Email:     newMainUser.Email,
		Password:  newMainUser.Password,
		CreatedAt: createdAt,
	}
	user, err := repo.userRepository.CreateUser(&newUser, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user : %w", err)
	}

	newAccount := entity.Account{
		ID:        idAccount,
		UserID:    user.ID,
		Owner:     newMainUser.FirstName + " " + newMainUser.LastName,
		Balance:   0,
		Currency:  "Rp",
		CreatedAt: time.Now(),
		IsActive: true,
	}
	account, err := repo.accountRepository.CreateAccount(&newAccount, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create account : %w", err)
	}

	newResponse := dto.UserResponse{
		ID:       account.ID,
		UserName: account.Owner,
		Email:    user.Email,
		Balance:  account.Balance,
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return nil, fmt.Errorf("failed to create user and account : %w", errCommit)
	}

	return &newResponse, nil
}

func (repo *transactionRepository) CreateMainTransfer(newMainTransfer *dto.TransferRequest, fromAccountID string, fromBalance int64, toBalance int64) (*dto.TransferResponse, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create user : %w", err)
	}

	idTransfer, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid : %w", err)
	}
	createdAt := time.Now()

	newTransfer := entity.Transfer{
		ID:            idTransfer.String(),
		FromAccountID: fromAccountID,
		ToAccountID:   newMainTransfer.ToAccountID,
		Amount:        newMainTransfer.Amount,
		CreatedAt:     createdAt,
	}
	transfer, err := repo.transferRepository.CreateTransfer(&newTransfer, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer : %w", err)
	}

	newFromAccount := dto.AccountRequest{
		AccountID: fromAccountID,
		Balance:   fromBalance,
	}
	updateFromAccount, err := repo.accountRepository.UpdateAccount(&newFromAccount, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance : %w", err)
	}

	newToAccount := dto.AccountRequest{
		AccountID: newMainTransfer.ToAccountID,
		Balance:   toBalance,
	}
	updateToAccount, err := repo.accountRepository.UpdateAccount(&newToAccount, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance : %w", err)
	}

	newResponse := dto.TransferResponse{
		ID:            transfer.ID,
		FromAccountID: updateFromAccount.AccountID,
		ToAccountID:   updateToAccount.AccountID,
		Amount:        transfer.Amount,
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return nil, fmt.Errorf("failed to create transfer : %w", errCommit)
	}

	return &newResponse, nil
}

func (repo *transactionRepository) CreateMainEntry(newMainEntry *dto.EntryRequest, accountID string, balance int64) (*dto.EntryResponse, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to create user : %w", err)
	}

	idEntry, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid : %w", err)
	}
	createdAt := time.Now()

	newEntry := entity.Entry{
		ID:        idEntry.String(),
		AccountID: accountID,
		Amount:    newMainEntry.Amount,
		CreatedAt: createdAt,
	}
	entry, err := repo.entryRepository.CreateEntry(&newEntry, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry : %w", err)
	}

	newAccount := dto.AccountRequest{
		AccountID: entry.AccountID,
		Balance:   balance,
	}
	updateAccount, err := repo.accountRepository.UpdateAccount(&newAccount, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance : %w", err)
	}

	newResponse := dto.EntryResponse{
		AccountID: entry.AccountID,
		Amount:    entry.Amount,
		Balance:   updateAccount.Balance,
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return nil, fmt.Errorf("failed to create entry : %w", errCommit)
	}

	return &newResponse, nil
}

func validate(err error, message string, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		fmt.Println(err, "transaction rollback")
	} else {
		fmt.Println(message)
	}
}