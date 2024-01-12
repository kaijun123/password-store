package controller

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"password_store/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignInController(c *gin.Context, db *gorm.DB) {
	// bind the request body to the struct RawCredentials
	var rc models.RawCredentials
	c.Bind(&rc)

	// print out the values in the struct
	username := rc.Username
	password := rc.Password
	fmt.Println(username)
	fmt.Println(password)

	// check if the username exists in the db
	var storedCredentials models.StoredCredentials
	if result := db.First(&storedCredentials, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"sign-in": "fail",
			"error":   "no such user",
		})
		return
	}

	// retrieve the salt for the username
	salt := storedCredentials.Salt
	oldHash := storedCredentials.Hash
	fmt.Println(storedCredentials.Username)
	fmt.Println(salt)
	fmt.Println("oldHash:", oldHash)

	// calculate new hash
	combinedString := password + strconv.Itoa(salt)
	h := sha256.New()
	h.Write([]byte(combinedString))
	newHash := h.Sum(nil)
	fmt.Println("newHash:", newHash)

	// compare the new hash and the old hash
	isEqual := bytes.Equal(oldHash, newHash)
	if isEqual {
		c.JSON(http.StatusOK, gin.H{
			"sign-in": "successful",
			"error":   "null",
		})
		return
	}

	// sends back a sample response
	c.JSON(http.StatusOK, gin.H{
		"sign-in": "fail",
		"error":   "incorrect password",
	})
}
