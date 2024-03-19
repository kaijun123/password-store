package controller

import (
	"fmt"
	"net/http"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context, sessionManager *kvStore.SessionManager) {

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
			sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware

			// Delete current session
			if err := sessionManager.DeleteSession(sessionId); err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"logout": "fail",
					"error":  "unable to logout",
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"logout": "successful",
					"error":  "null",
				})
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"refresh-cookie": "fail",
				"error":          "not authenticated",
			})
			return
		}
	}
}
