package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"password_store/internal/database"
	"password_store/internal/session"
	"password_store/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignInController(c *gin.Context, db *gorm.DB, sessionManager *session.SessionManager) {
	// bind the request body to the struct RawCredentials
	var rc database.RawCredentials
	c.Bind(&rc)
	// print out the values in the struct
	username := rc.Username
	password := rc.Password
	fmt.Println(username)
	fmt.Println(password)

	// check if the cookie is sent in the request, obtain cookie from request header;
	sessionId, err := c.Cookie(sessionManager.GetCookieName())
	if err != nil {
		// bad request
		if err != http.ErrNoCookie {
			c.JSON(http.StatusBadRequest, gin.H{
				"sign-in": "fail",
				"error":   "bad request",
			})
			return
		} else {

			// Case 1: Login in using username and password
			// check if the username exists in the db
			// if result.Error != nil -> you managed to find a record -> username exists in the db
			var storedCredentials database.StoredCredentials
			if result := db.First(&storedCredentials, "username = ?", username); result.Error != nil {
				c.JSON(http.StatusForbidden, gin.H{
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
				c.JSON(http.StatusForbidden, gin.H{
					"sign-in": "fail",
					"error":   "incorrect username and password",
				})
				return
			}
		}
	} else {

		// Case 2: Log in using cookie
		// retrieve session info from the session store; check if session has expired
		session, err := sessionManager.GetSession(sessionId)
		if err != nil { // no such session stored in the server
			c.JSON(http.StatusForbidden, gin.H{
				"sign-in": "fail",
				"error":   "invalid session",
			})
			return
		} else { // has a session stored in the server
			if session.IsExpired() { // session expired
				sessionManager.DeleteSession(sessionId) // ignore error as of now
				c.JSON(http.StatusForbidden, gin.H{
					"sign-in": "fail",
					"error":   "session expired",
				})
				return
			} else { // valid session
				c.JSON(http.StatusOK, gin.H{
					"sign-in": "successful",
					"error":   "null",
				})
			}
		}
	}
}
