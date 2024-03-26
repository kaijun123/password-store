package main

import (
	"fmt"
	auth "password_store/internal/controller/auth"
	txn "password_store/internal/controller/txn"

	"password_store/internal/kvStore"
	"password_store/internal/middleware"

	"github.com/gin-gonic/gin"

	"password_store/internal/database"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
)

func main() {
	fmt.Println("hello")

	// Connect to the database
	// Seems like the format is "postgres://{POSTGRES_USER}:{host}:{container_port}/{POSTGRES_DB}"
	// dbURL := "postgres://pg:pass@db:5432/password_store"

	redis := &kvStore.Redis{}
	redis.CreateClient()

	sessionManager := kvStore.NewSessionManager(redis)
	idempManager := kvStore.NewIdempotencyManager(redis)

	database := &database.Database{}
	database.Init()
	database.AutoMigrate()
	database.Seed()

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/sign-up", func(c *gin.Context) {
				auth.SignUpController(c, database.Db, sessionManager)
			})

			authGroup.POST("/sign-in", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				auth.SignInController(c, database.Db, sessionManager)
			})

			authGroup.POST("/logout", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				auth.LogoutHandler(c, sessionManager)
			})

			authGroup.POST("/refresh-cookie", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				auth.RefreshCookieHandler(c, sessionManager)
			})
		}

		txnGroup := v1.Group("/txn")
		{
			txnGroup.POST("/transfer", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				middleware.IdempMiddleware(c, *sessionManager, *idempManager)
			}, func(c *gin.Context) {
				txn.TransferHandler(c, sessionManager, idempManager, database.Db)
			})

			txnGroup.POST("/deposit", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				middleware.IdempMiddleware(c, *sessionManager, *idempManager)
			}, func(c *gin.Context) {
				txn.DepositHandler(c, sessionManager, idempManager, database.Db)
			})

			txnGroup.POST("/withdraw", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				middleware.IdempMiddleware(c, *sessionManager, *idempManager)
			}, func(c *gin.Context) {
				txn.WithdrawHandler(c, sessionManager, idempManager, database.Db)
			})

			txnGroup.POST("/fetch-balance", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				txn.FetchBalanceHandler(c, database.Db)
			})
		}

		// admin := v1.Group("/admin")
		// {
		// 	admin.POST("/delete-user", func(c *gin.Context) {
		// 		controller.DeleteUserController(c, database.Db)
		// 	})
		// }
	}

	router.Run()
}
