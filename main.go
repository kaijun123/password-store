package main

import (
	"fmt"
	"log"
	"password_store/internal/controller"
	"password_store/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("hello")

	// Connect to the database
	// Seems like the format is "postgres://{POSTGRES_USER}:{POSTGRES_PASSWORD}:{port}/{POSTGRES_DB}"
	dbURL := "postgres://pg:pass@db:5432/password_store"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	db.AutoMigrate(&models.StoredCredentials{})

	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	router.POST("/sign_in", func(c *gin.Context) {
		controller.SignInController(c, db)
	})

	router.POST("/sign_up", func(c *gin.Context) {
		controller.SignUpController(c, db)
	})
	router.Run()
}
