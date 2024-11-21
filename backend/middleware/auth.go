package middleware

import (
	"fmt"
	"mini-telegram/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies initData for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		initData := c.Param("auth")

		fmt.Println(initData)
		if initData == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "initData is missing"})
			return
		}

		// Verify initData
		data, valid := utils.VerifyInitData(initData)
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid initData"})
			return
		}

		fmt.Println("2222")
		// Set verified data in context for use in handlers
		c.Set("telegram_id", data["user_id"])
		c.Set("username", data["username"])

		c.Next()
	}
}
