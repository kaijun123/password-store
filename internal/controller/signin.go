package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"password_store/internal/database"
	"password_store/internal/kvStore"
	"password_store/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignInController(c *gin.Context, db *gorm.DB, sessionManager *kvStore.SessionManager) {
	// bind the request body to the struct RawCredentials
	var rc database.RawCredentials
	c.Bind(&rc)
	// print out the values in the struct
	username := rc.Username
	password := rc.Password
	fmt.Println(username)
	fmt.Println(password)

	authStatus, hasAuthStatus := c.Get("Auth Status")
	fmt.Println("Auth Status: ", authStatus)

	if !hasAuthStatus { // this is an error
		c.JSON(http.StatusBadRequest, gin.H{
			"sign-in": "fail",
			"error":   "bad request",
		})
		return
	} else {
		if authStatus == "Authenticated" { // Authenticated
			c.JSON(http.StatusBadRequest, gin.H{
				"sign-in": "success",
				"error":   "null",
			})
			return
		} else { // Not Authenticated: Log in using username and password
			var storedCredentials database.StoredCredentials
			if result := db.First(&storedCredentials, "username = ?", username); result.Error != nil {
				fmt.Println("Cannot find credentials in database")
				c.JSON(http.StatusForbidden, gin.H{
					"sign-in": "fail",
					"error":   "no account",
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
			combinedString := password + salt
			newHash := util.Hash([]byte(combinedString))
			fmt.Println("newHash:", newHash)

			// compare the new hash and the old hash
			isEqual := bytes.Equal(oldHash, newHash)
			if isEqual {
				sessionId, err := sessionManager.SetSession(username)
				if err != nil {
					panic("[Sign-Up Controller]: Unable to set cookie")
				}

				// set the sessionId as the value of the cookie
				SetCookieHandler(c, sessionManager.GetCookieName(), sessionId, sessionManager.GetExpiryDuration(), "/", "localhost", false, false)
				c.JSON(http.StatusOK, gin.H{
					"sign-in": "successful",
					"error":   "null",
				})
				return
			} else {
				fmt.Println("hash is not equal")
				c.JSON(http.StatusForbidden, gin.H{
					"sign-in": "fail",
					"error":   "invalid credentials",
				})
				return
			}
		}
	}
}
