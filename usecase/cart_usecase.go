package usecase

import (
	"errors"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/repository"
)

type CartUsecase interface {
	CreateCart(newCart *dto.CartRequest) (*dto.CartResponse, error)
}

type cartUsecase struct {
	cartRepo repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartUsecase(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartUsecase {
	return &cartUsecase{cartRepo: cartRepo, productRepo: productRepo}
}

func (usecase *cartUsecase) CreateCart(newCart *dto.CartRequest) (*dto.CartResponse, error) {
	if newCart.CustomerID == "" {
		return nil, errors.New("customer id is required")
	}
	if newCart.ProductID == "" {
		return nil, errors.New("product id is required")
	}
	if newCart.Quantity < 0 {
		return nil, errors.New("quantity is required")
	}
	getProductByProductID, err := usecase.productRepo.FindProductByProductID(newCart.ProductID)
	if err != nil {
		return nil, errors.New("product with product id " + newCart.ProductID + " not found")
	}
	stock := getProductByProductID.Stock - newCart.Quantity
	if stock < 0 {
		return nil, errors.New("stock not enough")
	}
	newCart.Price = getProductByProductID.Price
	total := getProductByProductID.Price * newCart.Quantity
	newCart.Total = total
	return usecase.cartRepo.CreateCart(newCart)
}