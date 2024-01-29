package main

import (
	"fmt"
	"password_store/internal/controller"
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
	router.POST("v1/sign-in", func(c *gin.Context) {
		controller.SignInController(c, database.Db, sessionManager)
	})

	router.POST("v1/sign-up", func(c *gin.Context) {
		controller.SignUpController(c, database.Db, sessionManager)
	})

	router.POST("v1/refresh-cookie", func(c *gin.Context) {
		controller.RefreshCookieHandler(c, sessionManager)
	})

	router.POST("v1/delete-user", func(c *gin.Context) {
		controller.DeleteUserController(c, database.Db)
	})

	router.Run()
}
