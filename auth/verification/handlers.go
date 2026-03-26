package verification

import (
	"log"
	"net/http"
	"os"

	"github.com/Astrasv/go-gully-backend/auth/local"
	"github.com/Astrasv/go-gully-backend/email"
	"github.com/gin-gonic/gin"
)

func VerifyEmailRedirect(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		frontendURL := os.Getenv("FRONTEND_URL")
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/verify-email?error=missing_token")
		return
	}

	_, err := VerifyUser(token)
	if err != nil {
		frontendURL := os.Getenv("FRONTEND_URL")
		if err == ErrUserNotFound || err == ErrTokenInvalid {
			c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/verify-email?error=invalid_token")
			return
		}
		if err == ErrUserAlreadyVerified {
			c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/verify-email?status=already_verified")
			return
		}
		log.Printf("[VERIFICATION] Verify error: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/verify-email?error=verification_failed")
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/verify-email?status=success")
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

func VerifyEmailAPI(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "token is required"})
		return
	}

	_, err := VerifyUser(req.Token)
	if err != nil {
		if err == ErrUserNotFound || err == ErrTokenInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_token", "message": "this verification link is invalid or has expired"})
			return
		}
		if err == ErrUserAlreadyVerified {
			c.JSON(http.StatusConflict, gin.H{"error": "already_verified", "message": "this email has already been verified"})
			return
		}
		log.Printf("[VERIFICATION] Verify error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "verification_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email_verified_successfully"})
}

func ResendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "valid email is required"})
		return
	}

	user, err := local.FindUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found", "message": "no account found with this email"})
		return
	}

	if user.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already_verified", "message": "this email has already been verified"})
		return
	}

	token, err := local.GenerateVerificationToken(user)
	if err != nil {
		log.Printf("[VERIFICATION] Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "resend_failed"})
		return
	}

	if err := email.SendVerificationEmail(user.Email, user.Username, token); err != nil {
		log.Printf("[VERIFICATION] Email send error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email_send_failed", "message": "failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification_email_sent"})
}
