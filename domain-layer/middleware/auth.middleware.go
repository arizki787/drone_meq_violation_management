package middleware

import (
	"BE-DroneMEQ-Violation-Management/domain-layer/service/external"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userManagementService external.UserManagementService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" || len(token) < 7 || token[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token
		token = token[7:]

		_, err := userManagementService.VerifyToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		userData, err := userManagementService.GetMetadataByToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to retrieve user data"})
			c.Abort()
			return
		}

		// Data context to use in controller
		c.Set("userData", userData)

		c.Next()
	}
}
