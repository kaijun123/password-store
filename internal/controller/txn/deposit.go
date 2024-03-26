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

// Used when a person deposits cash into the bank. Needs to increment the amount in person's bank account
func DepositHandler(c *gin.Context, sessionManager *kvStore.SessionManager, idempManager *kvStore.IdempotencyManager, db *gorm.DB) {
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
			if strings.ToLower(userTransaction.Type) != "deposit" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": constants.IdempBadRequest,
				})
				return
			}

			toPerson := userTransaction.To
			amt := userTransaction.Amt
			if toPerson == "" || amt == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": constants.IdempBadRequest,
				})
				return
			}

			// Check if the toPerson is the same as the user
			// ie the user is only authorized to deposit money into their own bank account
			sess, hasSess := c.Get("session")
			// an authorized user should always have the session set
			if !hasSess {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": constants.AuthServerErr,
				})
				return
			}

			session, hasSession := sess.(kvStore.Session)
			if !hasSession {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": constants.AuthServerErr,
				})
				return
			}

			if toPerson != session.Username {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid deposit account",
				})
				return
			}

			// The following methods are not thread-safe. Need some form of locking mechanism

			// Read the amount that the user currently has
			var balance database.UserBalance
			if err := db.Where("username = ?", toPerson).First(&balance).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to find user account",
				})
				return
			}

			// Write the new amount
			initAmt := balance.Balance
			newAmt := initAmt + amt
			balance.Balance = newAmt

			if err := db.Model(database.UserBalance{}).Where("username = ?", toPerson).Update("balance", newAmt).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to update the balance",
				})
				return
			}

			// Successful operation
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
			})

			// Write/ update the idemp status
			sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware
			userGenIdempKey := c.GetHeader(idempManager.GetRequestHeader())
			if idempStatus == constants.IdempFailed {
				// fmt.Println("updated status from failed to success")
				idempManager.UpdateIdempotency(userGenIdempKey, sessionId, constants.IdempSuccess, userTransactionBytes.([]byte))
			} else {
				// fmt.Println("wrote idemp status as success")
				idempManager.SetIdempotency(userGenIdempKey, sessionId, constants.IdempSuccess, userTransactionBytes.([]byte))
			}
		}
	}
}
