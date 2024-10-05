package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/adycahyoputro/merchant/model/dto"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("error loading .env file")
	}

	return os.Getenv(key)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CheckError(message string, err error) dto.Response {
	var result dto.Response
	if err.Error() == "internal server error" {
		result = dto.Response{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL_SERVER_ERROR",
			Data:   err.Error(),
		}
	}
	if err.Error() == "first name is required" || err.Error() == "last name is required" || err.Error() == "email is required" || err.Error() == "password is required" || err.Error() == "destination account is required" || err.Error() == "amount is required" || err.Error() == "amount must be positive amount" || err.Error() == "wrong password" || err.Error() == "store name is required" {
		result = dto.Response{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
	}
	if err.Error() == "user with email "+message+" has registered" {
		result = dto.Response{
			Code:   http.StatusFound,
			Status: "FOUND",
			Data:   err.Error(),
		}
	}
	if err.Error() == "balance is not enough" {
		result = dto.Response{
			Code:   http.StatusRequestedRangeNotSatisfiable,
			Status: "Range Not Satisfiable",
			Data:   err.Error(),
		}
	}
	if err.Error() == "account with account id "+message+" not found" || err.Error() == "user with email "+message+" not registered" {
		result = dto.Response{
			Code:   http.StatusNotFound,
			Status: "NOT_FOUND",
			Data:   err.Error(),
		}
	}
	if err.Error() == "user unauthorize" {
		result = dto.Response{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   err.Error(),
		}
	}
	return result
}

func AuthMiddleware(jwtKey string) gin.HandlerFunc {
	var jwtKeyByte = []byte(jwtKey)

	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtKeyByte, nil
		})

		if !token.Valid || err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		ctx.Set("claims", claims)
		ctx.Next()
	}
}