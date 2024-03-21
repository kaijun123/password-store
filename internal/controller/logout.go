package controller

import (
	"fmt"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context, sessionManager *kvStore.SessionManager) {
	authStatus, _ := c.Get(constants.AuthStatus)
	fmt.Println("Auth Status: ", authStatus)

	// No cookie
	if authStatus == constants.AuthNoCookie {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": authStatus,
		})
		return
	}

	// Authenticated
	if authStatus == constants.AuthAuthenticated {
		sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware

		// Delete current session
		if err := sessionManager.DeleteSession(sessionId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": constants.AuthServerErr,
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
			})
			return
		}
	}
}
