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
	assert.Error(t, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepository(db)

	user := &models.User{
		Email:    "user@example.com",
		Username: "user1",
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	retrievedUser, err := repo.GetByEmail("user@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, "user@example.com", retrievedUser.Email)
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepository(db)

	user := &models.User{
		Email:    "test2@example.com",
		Username: "user2",
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	retrievedUser, err := repo.GetUserByUsername("user2")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, "user2", retrievedUser.Username)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepository(db)

	user, err := repo.GetByEmail("notfound@example.com")
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserRepository_GetUserByUsername_NotFound(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepository(db)

	user, err := repo.GetUserByUsername("nonexistent")
	assert.Nil(t, user)
	assert.Error(t, err)
}

