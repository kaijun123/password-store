package middleware

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("entered here")
		fmt.Println("authStatus: ", authStatus)
		fmt.Println("hasAuthStatus: ", hasAuthStatus)

		c.Set(constants.IdempStatus, constants.IdempBadRequest)
		c.Next()
		return
	} else {
		// Get userGenIdempKey
		userGenIdempKey := c.GetHeader(idempManager.GetRequestHeader())

		// No userGenIdempKey
		if userGenIdempKey == "" {
			c.Set(constants.IdempStatus, constants.IdempNoUserGenKey)

		} else {
			var userTransaction database.UserTransaction
			if err := c.BindJSON(&userTransaction); err != nil {
				c.Set(constants.IdempStatus, constants.IdempServerErr)

			} else {
				userTransactionBytes, err := json.Marshal(&userTransaction)
				if err != nil {
					c.Set(constants.IdempStatus, constants.IdempServerErr)

				} else {
					idemp, err := idempManager.GetIdempotency(userGenIdempKey, sessionId, userTransactionBytes)
					fmt.Println("idemp:", idemp)

					if err != nil {
						if err.Error() == constants.IdempNew {
							// New request
							c.Set(constants.IdempStatus, constants.IdempNew)
							c.Set(constants.IdempBytes, userTransactionBytes)

						} else {
							// Invalid userGenIdempKey: no such key OR invalid session id OR invalid hash
							c.Set(constants.IdempStatus, err.Error())
						}

					} else {
						// Valid userGenIdempKey
						c.Set(constants.IdempStatus, idemp.Status)
						c.Set(constants.IdempBytes, userTransactionBytes)
					}
				}
			}
		}
	}

	// Does not handle IdempStatus = IdempNew and IdempStatus = IdempFailed
	// IdempStatus = IdempBadRequest is handled on top
	idempStatus, hasIdempStatus := c.Get(constants.IdempStatus)
	if !hasIdempStatus || idempStatus == constants.IdempServerErr {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.IdempServerErr,
		})
	} else if idempStatus == constants.IdempNoUserGenKey {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": idempStatus,
		})
	} else if idempStatus == constants.IdempSuccess {
		c.JSON(http.StatusOK, gin.H{
			"status": idempStatus,
		})
	}

	c.Next()
}
