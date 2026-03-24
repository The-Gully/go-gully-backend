package handlers

import (
	"github.com/gin-gonic/gin"
)

func Validate(c *gin.Context) {
	c.JSON(200, gin.H{"message": "You are authenticated"})
}

func Protected(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, gin.H{
		"message": "This is a secured endpoint",
		"user":    user,
	})
}
