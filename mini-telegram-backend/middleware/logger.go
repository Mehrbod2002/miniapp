package middleware

import (
	"time"

	"mini-telegram/config"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Continue to the next middleware or handler
		c.Next()

		// Perform logging asynchronously
		go func() {
			latency := time.Since(start)
			status := c.Writer.Status()
			config.Logger.WithFields(map[string]interface{}{
				"status":   status,
				"method":   c.Request.Method,
				"path":     c.Request.URL.Path,
				"latency":  latency.String(),
				"clientIP": c.ClientIP(),
			}).Info("Request completed")
		}()
	}
}
