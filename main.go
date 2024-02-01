package main

import (
	"fmt"
	"password_store/internal/controller"
	"password_store/internal/middleware"
	"password_store/internal/session"

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

	redis := &session.Redis{}
	redis.CreateClient()

	sessionManager := session.NewDefaultManager(redis)

	database := &database.Database{}
	database.Init()
	database.AutoMigrate()

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", func(c *gin.Context) {
				controller.SignUpController(c, database.Db, sessionManager)
			})

			auth.POST("/sign-in", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				controller.SignInController(c, database.Db, sessionManager)
			})

			auth.POST("/logout", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				controller.LogoutHandler(c, sessionManager)
			})

			auth.POST("/refresh-cookie", func(c *gin.Context) {
				middleware.AuthMiddleware(c, *sessionManager)
			}, func(c *gin.Context) {
				controller.RefreshCookieHandler(c, sessionManager)
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
