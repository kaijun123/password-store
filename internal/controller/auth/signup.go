package controller

import (
	"fmt"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/database"
	"password_store/internal/kvStore"
	"password_store/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignUpController(c *gin.Context, db *gorm.DB, sessionManager *kvStore.SessionManager) {
	// bind the request body to the struct RawUserCredentials
	var rc database.RawUserCredentials
	if err := c.BindJSON(&rc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.AuthBadRequest,
		})
		return
	}

	// print out the values in the struct
	username := rc.Username
	password := rc.Password

	fmt.Println(username)
	fmt.Println(password)

	// check if the username provided is used. Each username can only be used once
	// If result.Error == nil -> you didn't manage to find a record -> username is already in use
	var storedUserCredentials database.StoredUserCredentials
	if result := db.First(&storedUserCredentials, "username = ?", username); result.Error == nil {
		// fmt.Println(result.Error)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "username already used",
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

	storedUserCredentials = database.StoredUserCredentials{
		Username: username,
		Salt:     salt,
		Hash:     hash,
	}

	if result := db.Create(&storedUserCredentials); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.AuthServerErr,
		})
		return
	}

	// TODO: Create a row in the UserBalance database

	// sends back a sample response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
