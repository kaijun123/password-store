package controller

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"password_store/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignUpController(c *gin.Context, db *gorm.DB) {
	// bind the request body to the struct RawCredentials
	var rc models.RawCredentials
	c.Bind(&rc)

	// print out the values in the struct
	username := rc.Username
	password := rc.Password

	fmt.Println(username)
	fmt.Println(password)

	// check if the username provided is used. Each username can only be used once
	var storedCredentials models.StoredCredentials
	if result := db.First(&storedCredentials, "username = ?", username); result.Error == nil {
		// fmt.Println(result.Error)
		c.JSON(http.StatusOK, gin.H{
			"sign-up": "fail",
			"error":   "username already used",
		})
		return
	}

	// generate the salt
	salt := rand.Intn(10000000000000)
	fmt.Println(salt)

	combinedString := password + strconv.Itoa(salt)
	fmt.Println(combinedString)

	// compute the hash
	h := sha256.New()
	h.Write([]byte(combinedString))
	hash := h.Sum(nil)
	fmt.Printf("%x", hash)

	storedCredentials = models.StoredCredentials{
		Username: username,
		Salt:     salt,
		Hash:     hash,
	}

	if result := db.Create(&storedCredentials); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"sign-up": "fail",
			"error":   result.Error,
		})
	}

	// sends back a sample response
	c.JSON(http.StatusOK, gin.H{
		"sign-up": "successful",
		"error":   "null",
	})
}
