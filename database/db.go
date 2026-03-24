package database

import (
	"github.com/Astrasv/go-gully-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect(dsn string) *gorm.DB {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}
	return db
}

func Migrate() {
	db.AutoMigrate(&models.User{}, &models.Query{})
}

func GetDB() *gorm.DB {
	return db
}
