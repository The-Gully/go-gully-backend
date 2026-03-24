package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"unique"`
	Password   string
	ProviderID string `gorm:"index"`
	Role       string `gorm:"default:user"`
	Provider   string
	AvatarURL  string
	Name       string
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}
