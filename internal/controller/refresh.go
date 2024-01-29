package controller

import (
	"net/http"
	"password_store/internal/session"

	"github.com/gin-gonic/gin"
)

func RefreshCookieHandler(c *gin.Context, sessionManager *session.SessionManager) {
	sessionId, err := c.Cookie(sessionManager.GetCookieName())
	if err != nil {
		if err != http.ErrNoCookie { // bad request
			c.JSON(http.StatusBadRequest, gin.H{
				"refresh-cookie": "fail",
				"error":          "bad request",
			})
			return
		} else { // no cookie
			c.JSON(http.StatusForbidden, gin.H{
				"refresh-cookie": "fail",
				"error":          "no cookie",
			})
			return
		}
	}

	existingSession, err := sessionManager.GetSession(sessionId)
	if err != nil { // no such session stored in the server
		c.JSON(http.StatusForbidden, gin.H{
			"refresh-cookie": "fail",
			"error":          "invalid session",
		})
		return
	} else { // has a session stored in the server
		if existingSession.IsExpired() { // session expired
			c.JSON(http.StatusForbidden, gin.H{
				"refresh-cookie": "fail",
				"error":          "session expired",
			})
			return
		} else { // valid session: create new session and delete old session

			newSessionId, err := sessionManager.SetSession(existingSession.Username)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
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
				c.JSON(http.StatusForbidden, gin.H{
					"refresh-cookie": "successful",
					"error":          "null",
				})
				return
			}
		}
	}
}
