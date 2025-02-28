package repository

import (
	"auth-service/internal/models"
	"fmt"
    "gorm.io/gorm"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
)
type UserRepository struct {
	*repository.GenericRepository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{GenericRepository: repository.NewGenericRepository[models.User](db)}
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := repo.Db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to get user by email: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.Db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
