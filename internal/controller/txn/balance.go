package controller

import (
	"fmt"
	"net/http"
	"password_store/internal/constants"
	"password_store/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchBalanceHandler(c *gin.Context, db *gorm.DB) {
	authStatus, _ := c.Get(constants.AuthStatus)
	fmt.Println("Auth Status: ", authStatus)

	// No cookie
	if authStatus == constants.AuthNoCookie {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": authStatus,
		})
		return
	}

	// Authenticated
	if authStatus == constants.AuthAuthenticated {

		var userBalance database.UserBalance
		if err := c.BindJSON(&userBalance); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": constants.IdempServerErr,
			})
			return
		}

		if err := db.Where("username = ?", userBalance.Username).First(&userBalance).Error; err != nil {
			// User does not exist/ DB issue
			c.JSON(http.StatusBadRequest, gin.H{
				"error": constants.IdempBadRequest,
			})
		} else {
			// Return the balance
			c.JSON(http.StatusOK, gin.H{
				"username": userBalance.Username,
				"balance":  userBalance.Balance,
			})
		}
	}
}
