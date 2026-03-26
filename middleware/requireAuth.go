package middleware

import (
	"log"
	"net/http"

	"github.com/Astrasv/go-gully-backend/auth/google"
	"github.com/Astrasv/go-gully-backend/models"
	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	store := google.GetStore()
	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Printf("[AUTH] Middleware session error: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, exists := session.Values[google.SessionKey]
	if !exists || userID == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		log.Printf("[AUTH] Invalid user ID type in session")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := google.FindUserByID(id)
	if err != nil {
		log.Printf("[AUTH] User not found: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if user.Role == "" {
		user.Role = string(models.RoleUser)
	}

	c.Set("user", user)
	c.Set("userID", user.ID)
	c.Set("userRole", user.Role)
	c.Next()
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient_permissions"})
	}
}
