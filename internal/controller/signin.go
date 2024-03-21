package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/database"
	"password_store/internal/kvStore"
	"password_store/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignInController(c *gin.Context, db *gorm.DB, sessionManager *kvStore.SessionManager) {
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

	authStatus, _ := c.Get(constants.AuthStatus)
	fmt.Println("Auth Status: ", authStatus)

	// Authenticated
	if authStatus == constants.AuthAuthenticated {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
		return
	}

	// No cookie: Log in using username and password
	if authStatus == constants.AuthNoCookie {
		var storedUserCredentials database.StoredUserCredentials
		if result := db.First(&storedUserCredentials, "username = ?", username); result.Error != nil {
			fmt.Println("Cannot find credentials in database")
			c.JSON(http.StatusForbidden, gin.H{
				"error": constants.AuthInvalidCredentials,
			})
			return
		}

		// retrieve the salt for the username
		salt := storedUserCredentials.Salt
		oldHash := storedUserCredentials.Hash
		fmt.Println(storedUserCredentials.Username)
		fmt.Println(salt)
		fmt.Println("oldHash:", oldHash)

		// calculate new hash
		combinedString := password + salt
		newHash := util.Hash([]byte(combinedString))
		fmt.Println("newHash:", newHash)

		// compare the new hash and the old hash
		isEqual := bytes.Equal(oldHash, newHash)
		if isEqual {
			sessionId, err := sessionManager.SetSession(username)
			if err != nil {
				panic("[Sign-In Controller]: Unable to set cookie")
			}

			// set the sessionId as the value of the cookie
			SetCookieHandler(c, sessionManager.GetCookieName(), sessionId, sessionManager.GetExpiryDuration(), "/", "localhost", false, false)
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
			})
			return
		} else {
			fmt.Println("hash is not equal")
			c.JSON(http.StatusForbidden, gin.H{
				"error": constants.AuthInvalidCredentials,
			})
			return
		}
	}
}
