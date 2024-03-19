package controller

import (
	"fmt"
	"net/http"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
)

func RefreshCookieHandler(c *gin.Context, sessionManager *kvStore.SessionManager) {

	authStatus, hasAuthStatus := c.Get("Auth Status")
	fmt.Println("Auth Status: ", authStatus)

	if !hasAuthStatus { // this is an error
		c.JSON(http.StatusBadRequest, gin.H{
			"sign-in": "fail",
		})
		return
	} else {
		if authStatus == "Authenticated" { // Authenticated
			sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware
			existingSession, _ := sessionManager.GetSession(sessionId)

			newSessionId, err := sessionManager.SetSession(existingSession.Username)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"refresh-cookie": "fail",
					"error":          "unable to refresh cookie",
				})
				return
			}

			// Unable to delete old session. But it's ok as redis will remove the expired key-value pairs on its own
			// Just need to make sure that the new sessionId is not made available to the client
			if err := sessionManager.DeleteSession(sessionId); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"refresh-cookie": "fail",
					"error":          "unable to refresh cookie",
				})
				return
			} else {
				SetCookieHandler(c, sessionManager.GetCookieName(), newSessionId, sessionManager.GetExpiryDuration(), "/", "localhost", false, false)
				c.JSON(http.StatusOK, gin.H{
					"refresh-cookie": "successful",
					"error":          "null",
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
