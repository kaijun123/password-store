package middleware

import (
	"net/http"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
)

// TODO: Create an enum for the AuthStatus

// checks if there is an existing session
func AuthMiddleware(c *gin.Context, sessionManager kvStore.SessionManager) {
	sessionId, err := c.Cookie(sessionManager.GetCookieName())
	if err != nil {
		// bad request
		if err == http.ErrNoCookie {
			c.Set("Auth Status", "No Cookie")
		} else {
			c.Set("Auth Status", "Bad Request")
		}
		c.Next()
		return
	} else {
		session, err := sessionManager.GetSession(sessionId)
		if err != nil {
			c.Set("Auth Status", "No Session")
		} else {
			if session.IsExpired() {
				sessionManager.DeleteSession(sessionId) // ignore error as of now
				c.Set("Auth Status", "No Session")
			} else {
				c.Set("Auth Status", "Authenticated")
			}
		}
		c.Next()
		return
	}
}
