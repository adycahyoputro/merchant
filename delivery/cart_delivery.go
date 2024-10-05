package delivery

import (
	"net/http"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/usecase"
	"github.com/adycahyoputro/merchant/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CartDelivery interface {
	CreateCart(*gin.Context)
}

type cartDelivery struct {
	cartUsecase usecase.CartUsecase
}

func NewCartDelivery(
	cartUsecase usecase.CartUsecase) CartDelivery {
	return &cartDelivery{
		cartUsecase: cartUsecase,
	}
}

func (delivery *cartDelivery) CreateCart(ctx *gin.Context) {
	var cartRequest dto.CartRequest

	claims := ctx.MustGet("claims").(jwt.MapClaims)
	userId := claims["id"].(string)

	err := ctx.ShouldBindJSON(&cartRequest)
	if err != nil {
		result := dto.Response{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL_SERVER_ERROsR",
			Data:   err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	cartRequest.CustomerID = userId

	newCart, err := delivery.cartUsecase.CreateCart(&cartRequest)
	if err != nil {
		result := utils.CheckError(cartRequest.CustomerID, err)
		ctx.JSON(result.Code, result)
		return
	}

	result := dto.Response{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   newCart,
	}
	ctx.JSON(http.StatusCreated, result)
}