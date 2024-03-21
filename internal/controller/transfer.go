package controller

import (
	"fmt"
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

	if authStatus == constants.AuthAuthenticated && idempStatus != constants.IdempSuccess { // Authenticated
		// sessionId, _ := c.Cookie(sessionManager.GetCookieName()) // Error already handled in Auth middleware

		var userTransaction database.UserTransaction
		c.BindJSON(&userTransaction) // TODO: shift transfer struct into gin.Context
		// transferBytes, err := json.Marshal(transfer)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": constants.IdempServerErr,
		// 	})
		// 	return
		// }

		// idempotencyKey := c.GetHeader(idempManager.GetRequestHeader())
		// idemp, _ := idempManager.GetIdempotency(idempotencyKey, sessionId, util.Hash(transferBytes)) // TODO: shift idemp struct into gin.Context

		fromPerson := userTransaction.From
		toPerson := userTransaction.To
		amt := userTransaction.Amt

		fmt.Println("fromPerson: ", fromPerson)
		fmt.Println("toPerson: ", toPerson)
		fmt.Println("amt: ", amt)

		// db.Transaction(func(tx *gorm.DB) error {
		// 	// First find the From person, and the account balance
		// 	var fromBalance database.UserBalance
		// 	tx.Find()

		// })

		// Find the fromPerson, verify account balance > amt

		// Find the toPerson

	}
}
