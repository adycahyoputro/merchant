package delivery

import (
	"net/http"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/adycahyoputro/merchant/usecase"
	"github.com/adycahyoputro/merchant/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type StoreDelivery interface {
	CreateStore(*gin.Context)
}

type storeDelivery struct {
	storeUsecase usecase.StoreUsecase
}

func NewStoreDelivery(
	storeUsecase usecase.StoreUsecase) StoreDelivery {
	return &storeDelivery{
		storeUsecase: storeUsecase,
	}
}

func (delivery *storeDelivery) CreateStore(ctx *gin.Context) {
	var storeRequest dto.StoreRequest

	claims := ctx.MustGet("claims").(jwt.MapClaims)
	userId := claims["id"].(string)

	err := ctx.ShouldBindJSON(&storeRequest)
	if err != nil {
		result := dto.Response{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL_SERVER_ERROsR",
			Data:   err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	storeRequest.UserID = userId

	newStore, err := delivery.storeUsecase.CreateStore(&storeRequest)
	if err != nil {
		result := utils.CheckError(storeRequest.UserID, err)
		ctx.JSON(result.Code, result)
		return
	}

	result := dto.Response{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   newStore,
	}
	ctx.JSON(http.StatusCreated, result)
}