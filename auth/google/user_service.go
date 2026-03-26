package google

import (
	"github.com/Astrasv/go-gully-backend/database"
	"github.com/Astrasv/go-gully-backend/models"
)

func FindUserByID(id uint) (*models.User, error) {
	var user models.User
	result := database.GetDB().First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func FindOrCreateUser(userInfo *OAuth2UserInfo) (*models.User, error) {
	var user models.User

	result := database.GetDB().Where("email = ?", userInfo.Email).First(&user)
	if result.Error == nil {
		user.Name = userInfo.Name
		user.AvatarURL = userInfo.Picture
		if user.Provider == "" {
			user.Provider = "google"
			user.ProviderID = userInfo.ID
		}
		database.GetDB().Save(&user)
		return &user, nil
	}

	user = models.User{
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		AvatarURL:  userInfo.Picture,
		ProviderID: userInfo.ID,
		Provider:   "google",
		Role:       string(models.RoleUser),
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

type OAuth2UserInfo struct {
	ID      string
	Email   string
	Name    string
	Picture string
}
