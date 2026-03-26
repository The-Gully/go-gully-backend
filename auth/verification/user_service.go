package verification

import (
	"errors"
	"time"

	"github.com/Astrasv/go-gully-backend/database"
	"github.com/Astrasv/go-gully-backend/models"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrTokenInvalid        = errors.New("invalid or expired token")
	ErrUserAlreadyVerified = errors.New("user already verified")
)

func FindUserByVerificationToken(token string) (*models.User, error) {
	var user models.User
	result := database.GetDB().Where("verification_token = ?", token).First(&user)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func VerifyUser(token string) (*models.User, error) {
	user, err := FindUserByVerificationToken(token)
	if err != nil {
		return nil, err
	}

	if user.EmailVerified {
		return nil, ErrUserAlreadyVerified
	}

	if time.Now().After(user.VerificationExpiresAt) {
		return nil, ErrTokenInvalid
	}

	user.EmailVerified = true
	user.VerificationToken = ""
	user.VerificationExpiresAt = time.Time{}

	if err := database.GetDB().Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
