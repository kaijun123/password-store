package controller

import (
	"fmt"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
)

func RefreshCookieHandler(c *gin.Context, sessionManager *kvStore.SessionManager) {
	authStatus, _ := c.Get("Auth Status")
	fmt.Println("Auth Status: ", authStatus)

	// No cookie
	if authStatus == constants.AuthNoCookie {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": authStatus,
		})
		return
	}

	if authStatus == constants.AuthAuthenticated { // Authenticated
		sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware
		existingSession, _ := sessionManager.GetSession(sessionId)

		newSessionId, err := sessionManager.SetSession(existingSession.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": constants.AuthServerErr,
			})
			return
		}

		// Unable to delete old session. But it's ok as redis will remove the expired key-value pairs on its own
		// Just need to make sure that the new sessionId is not made available to the client
		if err := sessionManager.DeleteSession(sessionId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": constants.AuthServerErr,
			})
			return
		} else {
			SetCookieHandler(c, sessionManager.GetCookieName(), newSessionId, sessionManager.GetExpiryDuration(), "/", "localhost", false, false)
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
			})
			return
		}
	}
}
