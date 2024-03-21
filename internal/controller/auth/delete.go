package controller

// import (
// 	"fmt"
// 	"net/http"
// 	"password_store/internal/database"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// // TODO: Need to remove the session too
// func DeleteUserController(c *gin.Context, db *gorm.DB) {
// 	// bind the request body to the struct RawUserCredentials
// 	var rc database.RawUserCredentials
// 	c.Bind(&rc)

// 	// print out the values in the struct
// 	username := rc.Username
// 	password := rc.Password

// 	fmt.Println(username)
// 	fmt.Println(password)

// 	// check if the username exists in the db
// 	// if result.Error == nil -> you managed to find a record -> username exists in the db -> need to delete the user
// 	var storedCredentials database.StoredCredentials
// 	if result := db.First(&storedCredentials, "username = ?", username); result.Error == nil {
// 		db.Delete(&storedCredentials, "username = ?", username)
// 	}

// 	// sends back a sample response
// 	c.JSON(http.StatusOK, gin.H{
// 		"sign-up": "successful",
// 		"error":   "null",
// 	})
// }
