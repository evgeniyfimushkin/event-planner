package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB open connection with postgres db
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db
}

// SetupDB setup postgres and migrates all models
func SetupDB(dsn string, models ...interface{}) *gorm.DB {
	db := InitDB(dsn)
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Error migrating schema: %v", err)
	}
	return db
}

