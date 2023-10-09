package middleware

import (
	"net/http"

	"golang-api/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Token required",
			})
			c.Abort()
			return
		}

		token := string(tokenString[7:])

		if utils.RevokedTokens[token] {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Token revoked",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Invalid token",
			})
			c.Abort()
			return
		}
		username := claims["username"].(string)
		c.Set("claims", claims)
		c.Set("username", username)
		c.Next()
	}
}
