package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"user-service/pkg/service"
)

func JWTMiddleware(middlewareService *service.MiddlewareService) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		// Validate the token with the authentication service
		authResponse, err := middlewareService.ValidateTokenWithAuthService(c, token)
		if err != nil || !authResponse.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Store the authenticated user info in the context (optional)
		c.Set("user", authResponse.User)

		// Proceed to the next handler
		c.Next()

	}
}
