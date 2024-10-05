package app

import (
	"database/sql"

	"github.com/adycahyoputro/merchant/delivery"
	"github.com/adycahyoputro/merchant/repository"
	"github.com/adycahyoputro/merchant/usecase"
	"github.com/adycahyoputro/merchant/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	userRepository := repository.NewUserRepository(db)
	accountRepository := repository.NewAccountRepository(db)
	entryRepository := repository.NewEntryRepository(db)
	storeRepository := repository.NewStoreRepository(db)
	productRepository := repository.NewProductRepository(db)
	cartRepository := repository.NewCartRepository(db)
	transferRepository := repository.NewTransferRepository(db)
	transactionRepository := repository.NewTransactionRepository(db, userRepository, accountRepository, entryRepository, transferRepository)

	// userUsecase := usecase.NewUserUsecase(userRepository)
	accountUsecase := usecase.NewAccountUsecase(accountRepository)
	storeUsecase := usecase.NewStoreUsecase(storeRepository)
	productUsecase := usecase.NewProductUsecase(productRepository)
	cartUsecase := usecase.NewCartUsecase(cartRepository, productRepository)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepository, userRepository, accountRepository)
	loginUsecase := usecase.NewLoginUsecase(userRepository, accountRepository)

	storeDelivery := delivery.NewStoreDelivery(storeUsecase)
	productDelivery := delivery.NewProductDelivery(productUsecase)
	cartDelivery := delivery.NewCartDelivery(cartUsecase)
	transactionDelivery := delivery.NewTransactionDelivery(transactionUsecase, accountUsecase)
	loginDelivery := delivery.NewLoginDelivery(loginUsecase)

	r.POST("/login", loginDelivery.Login)
	r.POST("/register", transactionDelivery.CreateUserAccount)

	group_router := r.Group("merchant")
	group_router.Use(utils.AuthMiddleware(utils.GetEnv("JWT_KEY")))
	{
		group_router.POST("/entry", transactionDelivery.CreateMainEntry)
		group_router.POST("/transfer", transactionDelivery.CreateMainTransfer)
		group_router.POST("/logout", loginDelivery.Logout)
		group_router.POST("/store", storeDelivery.CreateStore)
		group_router.POST("/product", productDelivery.CreateProduct)
		group_router.POST("/cart", cartDelivery.CreateCart)
	}

	return r
}