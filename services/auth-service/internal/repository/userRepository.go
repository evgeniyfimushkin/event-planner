package repository

import (
	"auth-service/internal/models"
	"fmt"
    "gorm.io/gorm"
)
type UserRepository struct {
	*GenericRepository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{GenericRepository: NewGenericRepository[models.User](db)}
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to get user by email: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
