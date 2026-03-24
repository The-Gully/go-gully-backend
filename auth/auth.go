package auth

import (
	"os"

	"github.com/Astrasv/go-gully-backend/database"
	"github.com/Astrasv/go-gully-backend/models"
)

type User = models.User
type Role = models.Role

const (
	RoleUser  = models.RoleUser
	RoleAdmin = models.RoleAdmin
)

func Initialize(clientID, clientSecret, callbackURL, sessionSecret string) {
	InitOAuth(clientID, clientSecret, callbackURL)
	InitSession(sessionSecret)
}

func LoadEnvAndConnect() {
	database.Connect(os.Getenv("DB"))
	database.Migrate()
}

func GetUserByID(id uint) (*models.User, error) {
	return FindUserByID(id)
}
