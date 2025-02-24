package repository

import (
	"auth-service/internal/models"
	"fmt"
    "log"
    "gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error migrating schema: %v", err)
	}

	return &UserRepository{db: db}
}

func (repo *UserRepository) Create(user *models.User) error {
	result := repo.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("unable to create user: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (repo *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	result := repo.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil 
		}
		return nil, fmt.Errorf("unable to get user by id: %w", result.Error)
	}
	return &user, nil
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil 
		}
		return nil, fmt.Errorf("unable to get user by email: %w", result.Error)
	}
	return &user, nil
}

func (repo *UserRepository) Update(user *models.User) error {
	result := repo.db.Save(user)
	if result.Error != nil {
		return fmt.Errorf("unable to update user: %w", result.Error)
	}
	return nil
}

func (repo *UserRepository) Delete(id int) error {
	result := repo.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("unable to delete user: %w", result.Error)
	}
	return nil
}

func (repo *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	result := repo.db.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to fetch users: %w", result.Error)
	}
	return users, nil
}
