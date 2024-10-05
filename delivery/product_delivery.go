package delivery

import (
	"net/http"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/usecase"
	"github.com/adycahyoputro/merchant/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ProductDelivery interface {
	CreateProduct(*gin.Context)
}

type productDelivery struct {
	productUsecase usecase.ProductUsecase
}

func NewProductDelivery(
	productUsecase usecase.ProductUsecase) ProductDelivery {
	return &productDelivery{
		productUsecase: productUsecase,
	}
}

func (delivery *productDelivery) CreateProduct(ctx *gin.Context) {
	var productRequest dto.ProductRequest

	claims := ctx.MustGet("claims").(jwt.MapClaims)
	userId := claims["id"].(string)

	err := ctx.ShouldBindJSON(&productRequest)
	if err != nil {
		result := dto.Response{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL_SERVER_ERROsR",
			Data:   err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	// productRequest.StoreID = userId

	newProduct, err := delivery.productUsecase.CreateProduct(&productRequest, userId)
	if err != nil {
		result := utils.CheckError(productRequest.ProductName, err)
		ctx.JSON(result.Code, result)
		return
	}

	result := dto.Response{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   newProduct,
	}
	ctx.JSON(http.StatusCreated, result)
}