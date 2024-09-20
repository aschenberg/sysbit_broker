package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Role(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate getting user from context or other means
		user, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})

			return
		}

		// Ensure userRoles is of type []string
		actualRoles, ok := user.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user roles data"})
			return
		}

		// Check if the user has at least one of the required roles
		if !hasAnyRole(actualRoles, requiredRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access forbidden for role"})

			return
		}

		// User is authorized, proceed to next handler
		c.Next()
	}
}

func hasAnyRole(user []string, requiredRoles []string) bool {
	roleSet := make(map[string]struct{})
	for _, r := range user {
		roleSet[r] = struct{}{}
	}
	for _, r := range requiredRoles {
		if _, exists := roleSet[r]; exists {
			return true
		}
	}
	return false
}
