package models

import "time"

type Query struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"`
	Query     string `gorm:"type:text;not null"`
	Response  string `gorm:"type:text;not null"`
	CreatedAt time.Time
}
