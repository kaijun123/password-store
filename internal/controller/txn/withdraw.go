package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/database"
	"password_store/internal/kvStore"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WithdrawHandler(c *gin.Context, sessionManager *kvStore.SessionManager, idempManager *kvStore.IdempotencyManager, db *gorm.DB) {
	authStatus, _ := c.Get(constants.AuthStatus)
	fmt.Println("Auth Status: ", authStatus)

	idempStatus, _ := c.Get(constants.IdempStatus)
	fmt.Println("Idemp Status: ", idempStatus)

	// No cookie
	if authStatus == constants.AuthNoCookie {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": authStatus,
		})
		return
	}

	if authStatus == constants.AuthAuthenticated {
		if idempStatus == constants.IdempNew || idempStatus == constants.IdempFailed {
			userTransactionBytes, _ := c.Get(constants.IdempBytes)

			var userTransaction database.UserTransaction
			if err := json.Unmarshal(userTransactionBytes.([]byte), &userTransaction); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": constants.IdempServerErr,
				})
				return
			}

			// Check transaction type
			if strings.ToLower(userTransaction.Type) != "withdraw" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": constants.IdempBadRequest,
				})
				return
			}

		}
	}
}
