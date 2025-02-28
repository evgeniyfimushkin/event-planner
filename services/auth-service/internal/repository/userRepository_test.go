package repository

import (
	"auth-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	db.AutoMigrate(&models.User{})
	return db
}

func TestUserRepository_CRUD(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepository(db)

	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
	}

	err := repo.Create(user)
	assert.NoError(t, err)

	retrievedUser, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Username, retrievedUser.Username)

	user.Username = "updateduser"
	err = repo.Update(user)
	assert.NoError(t, err)

	updatedUser, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUser.Username)

	err = repo.Delete(user.ID)
	assert.NoError(t, err)

	deletedUser, err := repo.GetByID(user.ID)
	assert.Nil(t, deletedUser)
	assert.Nil(t, err)
}

