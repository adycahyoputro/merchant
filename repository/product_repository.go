package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/model/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	CreateProduct(newProduct *dto.ProductRequest) (*dto.ProductResponse, error)
	FindProductByProductID(productID string) (*entity.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (repo *productRepository) CreateProduct(newProduct *dto.ProductRequest) (*dto.ProductResponse, error) {
	stmt, err := repo.db.Prepare("INSERT INTO products (id, store_id, product_name, description, stock, price, created_at, updated_at, is_delete) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("failed to create product : %w", err)
	}
	defer stmt.Close()

	repeat := true
	var idProduct string
	for repeat {
		idProductUID, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to create uuid : %w", err)
		}
		if idProductUID.String() != newProduct.StoreID {
			repeat = false
			idProduct = idProductUID.String()
		}
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	isDelete := false
	err = stmt.QueryRow(idProduct, newProduct.StoreID, newProduct.ProductName, newProduct.Description, newProduct.Stock, newProduct.Price, createdAt, updatedAt, isDelete).Scan(idProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to create product : %w", err)
	}

	newResponse := dto.ProductResponse{
		ID: idProduct,
		StoreID: newProduct.StoreID,
		ProductName: newProduct.ProductName,
		Description: newProduct.Description,
		Stock: newProduct.Stock,
		Price: newProduct.Price,
	}
	return &newResponse, nil
}

func (repo *productRepository) FindProductByProductID(productID string) (*entity.Product, error) {
	var product entity.Product
	stmt, err := repo.db.Prepare("SELECT id, store_id, product_name, description, stock, price, created_at, updated_at, is_delete FROM products WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(productID).Scan(&product.ID, &product.StoreID, &product.ProductName, &product.Description, &product.Stock, &product.Price, &product.CreatedAt, &product.UpdateAt, &product.IsDelete)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product with user id %s not found", productID)
	} else if err != nil {
		return nil, err
	}

	return &product, nil
}