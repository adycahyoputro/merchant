package usecase

import (
	"errors"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/repository"
)

type ProductUsecase interface {
	CreateProduct(newProduct *dto.ProductRequest, userId string) (*dto.ProductResponse, error)
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (usecase *productUsecase) CreateProduct(newProduct *dto.ProductRequest, userId string) (*dto.ProductResponse, error) {
	if newProduct.ProductName == "" {
		return nil, errors.New("product name is required")
	}
	if newProduct.Description == "" {
		return nil, errors.New("description is required")
	}
	if newProduct.Price == 0 || newProduct.Price < 0{
		return nil, errors.New("price is required")
	}
	if newProduct.Stock < 0 {
		return nil, errors.New("stock is required")
	}
	return usecase.productRepo.CreateProduct(newProduct)
}