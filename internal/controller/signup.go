package controller

import (
	"fmt"
	"net/http"
	"password_store/internal/database"
	"password_store/internal/session"
	"password_store/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignUpController(c *gin.Context, db *gorm.DB, sessionManager *session.SessionManager) {
	// bind the request body to the struct RawCredentials
	var rc database.RawCredentials
	c.Bind(&rc)

	// print out the values in the struct
	username := rc.Username
	password := rc.Password

	fmt.Println(username)
	fmt.Println(password)

	// check if the username provided is used. Each username can only be used once
	// If result.Error == nil -> you didn't manage to find a record -> username is already in use
	var storedCredentials database.StoredCredentials
	if result := db.First(&storedCredentials, "username = ?", username); result.Error == nil {
		// fmt.Println(result.Error)
		c.JSON(http.StatusOK, gin.H{
			"sign-up": "fail",
			"error":   "username already used",
		})
		return
	}

	// generate the salt
	salt := util.GenerateRandomString(20)
	fmt.Println(salt)

	combinedString := password + salt
	fmt.Println(combinedString)

	// compute the hash
	hash := util.Hash([]byte(combinedString))
	fmt.Printf("%x", string(hash))

	storedCredentials = database.StoredCredentials{
		Username: username,
		Salt:     salt,
		Hash:     hash,
	}

	if result := db.Create(&storedCredentials); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"sign-up": "fail",
			"error":   result.Error,
		})
		return
	}

	// sends back a sample response
	c.JSON(http.StatusOK, gin.H{
		"sign-up": "successful",
		"error":   "null",
	})
}
