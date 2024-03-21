package middleware

import (
	"encoding/json"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/database"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
)

// checks if there is an existing session
func IdempMiddleware(c *gin.Context, sessionManager kvStore.SessionManager, idempManager kvStore.IdempotencyManager) {
	sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware

	// check if the user is authenticated
	authStatus, hasAuthStatus := c.Get(constants.AuthStatus)
	if !hasAuthStatus || authStatus != constants.AuthAuthenticated {
		c.Set(constants.IdempStatus, constants.IdempBadRequest)
	} else {
		// Create idempKey
		idempotencyKey := c.GetHeader(idempManager.GetRequestHeader())

		// No idempKey: New request
		if idempotencyKey == "" {
			c.Set(constants.IdempStatus, constants.IdempNew)
		} else {
			var UserTransaction database.UserTransaction
			if err := c.BindJSON(&UserTransaction); err != nil {
				c.Set(constants.IdempStatus, constants.IdempServerErr)

			} else {
				UserTransactionBytes, err := json.Marshal(&UserTransaction)
				if err != nil {
					c.Set(constants.IdempStatus, constants.IdempServerErr)

				} else {
					idemp, err := idempManager.GetIdempotency(idempotencyKey, sessionId, UserTransactionBytes)

					if err != nil {
						// Invalid idempKey: no such key OR invalid session id OR invalid hash
						c.Set(constants.IdempStatus, constants.IdempBadRequest)
					} else {
						// Valid idempKey
						c.Set(constants.IdempStatus, idemp.Status)
					}
				}
			}
		}
	}

	idempStatus, hasIdempStatus := c.Get(constants.IdempStatus)

	if !hasIdempStatus {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.IdempServerErr,
		})
	} else if idempStatus == constants.IdempSuccess {
		c.JSON(http.StatusOK, gin.H{
			"status": idempStatus,
		})
	}
	c.Next()
}
