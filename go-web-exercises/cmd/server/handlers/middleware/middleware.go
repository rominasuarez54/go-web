package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	tokenFromEnv := os.Getenv("TOKEN")
	return func(c *gin.Context) {
		tokenFromHeader := c.GetHeader("Token")

		if tokenFromEnv != tokenFromHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			return
		}
		c.Next()
	}
}
