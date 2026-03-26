package local

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/Astrasv/go-gully-backend/database"
	"github.com/Astrasv/go-gully-backend/models"
)

const TokenExpiryDuration = 24 * time.Hour

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
)

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.GetDB().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := database.GetDB().Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func FindUserByID(id uint) (*models.User, error) {
	var user models.User
	result := database.GetDB().First(&user, id)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func CreateUser(email, username, password string) (*models.User, error) {
	var existingUser models.User

	if err := database.GetDB().Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, ErrEmailExists
	}

	if err := database.GetDB().Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Provider: "local",
		Role:     string(models.RoleUser),
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func AuthenticateUser(emailOrUsername, password string) (*models.User, error) {
	var user models.User

	result := database.GetDB().Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user)
	if result.Error != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Provider != "local" {
		return nil, ErrInvalidCredentials
	}

	if !CheckPassword(password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

func GenerateVerificationToken(user *models.User) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)
	user.VerificationToken = token
	user.VerificationExpiresAt = time.Now().Add(TokenExpiryDuration)

	if err := database.GetDB().Save(user).Error; err != nil {
		return "", err
	}

	return token, nil
}
