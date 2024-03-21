package middleware

import (
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
)

// checks if there is an existing session
func AuthMiddleware(c *gin.Context, sessionManager kvStore.SessionManager) {
	sessionId, err := c.Cookie(sessionManager.GetCookieName())

	if err != nil {
		// bad request
		if err == http.ErrNoCookie {
			c.Set(constants.AuthStatus, constants.AuthNoCookie)
		} else {
			c.Set(constants.AuthStatus, constants.AuthBadRequest)
		}
	} else {
		session, err := sessionManager.GetSession(sessionId)
		if err != nil {
			c.Set(constants.AuthStatus, constants.AuthNoSession)
		} else {
			if session.IsExpired() {
				if err := sessionManager.DeleteSession(sessionId); err != nil {
					c.Set(constants.AuthStatus, constants.AuthServerErr)
				} else {
					c.Set(constants.AuthStatus, constants.AuthNoSession)
				}
			} else {
				c.Set(constants.AuthStatus, constants.AuthAuthenticated)
			}
		}
	}

	authStatus, hasAuthStatus := c.Get(constants.AuthStatus)

	// Does not handle AuthStatus = AuthNoCookie
	// This is because sign-in does not require a cookie
	// Individual Auth handlers will need to handle this error
	if !hasAuthStatus || authStatus != constants.AuthAuthenticated {
		if !hasAuthStatus || authStatus == constants.AuthServerErr {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": constants.AuthServerErr,
			})
		} else if authStatus == constants.AuthBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": authStatus,
			})
		} else if authStatus == constants.AuthNoSession {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": authStatus,
			})
		}
	}
	c.Next()
}
