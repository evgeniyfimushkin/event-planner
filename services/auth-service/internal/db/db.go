package db

import (
	"auth-service/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db
}

func SetupDB(dsn string) *gorm.DB {
	db := InitDB(dsn)
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error migrating schema: %v", err)
	}
	return db
}

