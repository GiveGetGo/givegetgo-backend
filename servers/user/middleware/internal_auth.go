package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// InternalAuthMiddleware - middleware to authenticate internal requests
func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get which service is calling
		service := c.GetHeader("X-Service")

		// Construct the environment variable name and retrieve the API key
		envVarName := service + "_API_KEY"
		expectedApiKey := os.Getenv(envVarName)

		// Check API key
		apiKey := c.GetHeader("X-Api-Key")
		if apiKey != expectedApiKey {
			c.JSON(403, gin.H{
				"code":    "40301",
				"message": "Forbidden - Invalid API Key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
