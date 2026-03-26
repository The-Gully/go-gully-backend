package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                 string    `gorm:"unique" json:"email"`
	Username              string    `gorm:"unique" json:"username"`
	Password              string    `json:"-"`
	ProviderID            string    `gorm:"index" json:"provider_id"`
	Role                  string    `gorm:"default:user" json:"role"`
	Provider              string    `json:"provider"`
	AvatarURL             string    `json:"avatar_url"`
	Name                  string    `json:"name"`
	EmailVerified         bool      `gorm:"default:false" json:"email_verified"`
	VerificationToken     string    `gorm:"index" json:"-"`
	VerificationExpiresAt time.Time `json:"-"`
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}
