package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/database"
	"password_store/internal/kvStore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransferHandler(c *gin.Context, sessionManager *kvStore.SessionManager, idempManager *kvStore.IdempotencyManager, db *gorm.DB) {

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

	// Authenticated
	if authStatus == constants.AuthAuthenticated {
		if idempStatus == constants.IdempNew || idempStatus == constants.IdempFailed {

			sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware
			userGenIdempKey := c.GetHeader(idempManager.GetRequestHeader())

			userTransactionBytes, _ := c.Get(constants.IdempBytes)

			var userTransaction database.UserTransaction
			if err := json.Unmarshal(userTransactionBytes.([]byte), &userTransaction); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": constants.IdempServerErr,
				})
				return
			}
			fmt.Println("userTransaction:", userTransaction)

			if userTransaction.Type != "transfer" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": constants.IdempBadRequest,
				})
				return
			}

			fromPerson := userTransaction.From
			toPerson := userTransaction.To
			amt := userTransaction.Amt

			fmt.Println("fromPerson: ", fromPerson)
			fmt.Println("toPerson: ", toPerson)
			fmt.Println("amt: ", amt)

			// TODO: Not thread-safe. Implement row-locking mechanism.
			err := db.Transaction(func(tx *gorm.DB) error {
				// First find the From person, and the account balance
				var fromBalance database.UserBalance
				if err := tx.Where("username = ?", fromPerson).First(&fromBalance).Error; err != nil {
					return errors.New("cannot get from balance")
				}
				fmt.Println("fromBalance", fromBalance)

				var toBalance database.UserBalance
				if err := tx.Where("username = ?", toPerson).First(&toBalance).Error; err != nil {
					return errors.New("cannot get to balance")
				}
				fmt.Println("toBalance", toBalance)

				// verify that fromBalance > amt
				fmt.Println("fromBalance.Balance - amt", fromBalance.Balance-amt)
				if (fromBalance.Balance - amt) < 0 {
					return errors.New("insufficient balance")
				}

				if err := tx.Model(&database.UserBalance{}).Where("username = ?", fromPerson).Update("balance", fromBalance.Balance-amt).Error; err != nil {
					return errors.New("cannot update from balance")
				}
				fmt.Println(toBalance)

				if err := tx.Model(&database.UserBalance{}).Where("username = ?", toPerson).Update("balance", toBalance.Balance+amt).Error; err != nil {
					return errors.New("cannot update to balance")
				}

				if err := tx.Create(&userTransaction).Error; err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				// Failed operation
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			} else {
				// Successful operation
				c.JSON(http.StatusOK, gin.H{
					"status": "success",
				})

				// Write/ update the idemp status
				if idempStatus == constants.IdempFailed {
					fmt.Println("updated status from failed to success")
					idempManager.UpdateIdempotency(userGenIdempKey, sessionId, constants.IdempSuccess, userTransactionBytes.([]byte))
				} else {
					fmt.Println("wrote idemp status as success")
					idempManager.SetIdempotency(userGenIdempKey, sessionId, constants.IdempSuccess, userTransactionBytes.([]byte))
				}
			}
		}
	}
}
