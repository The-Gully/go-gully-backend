package auth

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Login(c *gin.Context) {
	state := GenerateRandomString(32)

	session, err := GetStore().Get(c.Request, SessionName)
	if err != nil {
		log.Printf("[AUTH] Session error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session_error"})
		return
	}

	session.Values["oauth_state"] = state
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("[AUTH] Session save error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session_error"})
		return
	}

	url := GetOAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func Callback(c *gin.Context) {
	session, err := GetStore().Get(c.Request, SessionName)
	if err != nil {
		log.Printf("[AUTH] Callback session error: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=session_error")
		return
	}

	state := c.Query("state")
	storedState, ok := session.Values["oauth_state"].(string)
	if !ok || state != storedState {
		log.Printf("[AUTH] State mismatch")
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=invalid_state")
		return
	}

	code := c.Query("code")
	token, err := GetOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Printf("[AUTH] Token exchange error: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=token_exchange_failed")
		return
	}

	userInfo, err := FetchUserInfo(token.AccessToken)
	if err != nil {
		log.Printf("[AUTH] Get user info error: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=user_info_failed")
		return
	}

	user, err := FindOrCreateUser(&OAuth2UserInfo{
		ID:      userInfo.Id,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
	})
	if err != nil {
		log.Printf("[AUTH] Find or create user error: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=user_creation_failed")
		return
	}

	session.Values[SessionKey] = user.ID
	delete(session.Values, "oauth_state")

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("[AUTH] Session save error: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=session_save_failed")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"/dashboard")
}

func Logout(c *gin.Context) {
	session, err := GetStore().Get(c.Request, SessionName)
	if err != nil {
		log.Printf("[AUTH] Logout session error: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "logged_out"})
		return
	}

	session.Values[SessionKey] = nil
	session.Options.MaxAge = -1

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("[AUTH] Logout save error: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged_out"})
}

func GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
