package local

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Astrasv/go-gully-backend/auth/google"
	"github.com/Astrasv/go-gully-backend/models"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User     *models.User `json:"user"`
	Provider string       `json:"provider"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}

	if !isValidUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_username", "message": "username can only contain letters, numbers, and underscores"})
		return
	}

	user, err := CreateUser(req.Email, req.Username, req.Password)
	if err != nil {
		if err == ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "email_exists", "message": "an account with this email already exists"})
			return
		}
		if err == ErrUsernameExists {
			c.JSON(http.StatusConflict, gin.H{"error": "username_exists", "message": "this username is already taken"})
			return
		}
		log.Printf("[LOCAL AUTH] Registration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration_failed"})
		return
	}

	session, err := google.GetStore().Get(c.Request, google.SessionName)
	if err != nil {
		log.Printf("[LOCAL AUTH] Session error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session_error"})
		return
	}

	session.Values[google.SessionKey] = user.ID
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("[LOCAL AUTH] Session save error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session_error"})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	c.JSON(http.StatusCreated, gin.H{
		"message":  "registration_successful",
		"redirect": frontendURL + "/dashboard",
		"user":     user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "details": err.Error()})
		return
	}

	user, err := AuthenticateUser(req.EmailOrUsername, req.Password)
	if err != nil {
		if err == ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials", "message": "invalid email/username or password"})
			return
		}
		log.Printf("[LOCAL AUTH] Login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login_failed"})
		return
	}

	session, err := google.GetStore().Get(c.Request, google.SessionName)
	if err != nil {
		log.Printf("[LOCAL AUTH] Session error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session_error"})
		return
	}

	session.Values[google.SessionKey] = user.ID
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("[LOCAL AUTH] Session save error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session_error"})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	c.JSON(http.StatusOK, gin.H{
		"message":  "login_successful",
		"redirect": frontendURL + "/dashboard",
		"user":     user,
	})
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func isValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}
