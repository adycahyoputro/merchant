package repository

import (
	"database/sql"
	"fmt"

	"github.com/adycahyoputro/merchant/model/entity"
)

type UserRepository interface {
	FindUserByEmail(email string) (*entity.User, error)
	// FindUserBy(email string) (*entity.User, error)
	CreateUser(newUser *entity.User, tx *sql.Tx) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(
	db *sql.DB) UserRepository {
	return &userRepository{
		db: db}
}

func (repo *userRepository) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	stmt, err := repo.db.Prepare("SELECT id, first_name, last_name, email, password, created_at FROM users WHERE email = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with email %s not found", email)
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) CreateUser(newUser *entity.User, tx *sql.Tx) (*entity.User, error) {
	stmt, err := repo.db.Prepare("INSERT INTO users (id, first_name, last_name, email, password, created_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("failed to create user : %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(newUser.ID, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password, newUser.CreatedAt).Scan(&newUser.ID)

	validate(err, "create user", tx)

	return newUser, nil
}