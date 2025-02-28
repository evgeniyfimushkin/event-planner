package db

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string
}

func TestInitDB(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "should successfully open a connection to the database")
	assert.NotNil(t, database, "should return a database instance")
}

func TestSetupDB(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "should successfully open a connection to the database")
	assert.NotNil(t, database, "should return a database instance")

	err = database.AutoMigrate(&TestModel{})
	assert.NoError(t, err, "migration should complete without errors")

	var tableNames []string
	database.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name='test_models'").Scan(&tableNames)
	assert.Contains(t, tableNames, "test_models", "test_models table should exist in the database")
}

