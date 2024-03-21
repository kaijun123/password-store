package controller

import (
	"github.com/gin-gonic/gin"
)

// Reference: https://medium.com/@gopal96685/handling-cookies-with-gin-framework-in-go-f119358c9cf3#:~:text=Setting%20a%20Cookie%20in%20Gin&text=You%20can%20use%20the%20SetCookie,The%20value%20of%20the%20cookie.
func SetCookieHandler(c *gin.Context, name, value string, maxAge int, path, domain string, secure bool, httpOnly bool) {
	c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}
