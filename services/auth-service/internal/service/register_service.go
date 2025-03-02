package service

import (
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"fmt"
)

type RegisterService struct {
    userRepo *repository.UserRepository
}

func NewRegisterService(userRepo *repository.UserRepository) (*RegisterService, error) {
    return &RegisterService{
        userRepo: userRepo,
    }, nil
}

func (s *RegisterService) Register(username, email, passhash string) (*models.User, error) {
    if username == ""{
        return  nil, fmt.Errorf("username is required")
    }
    if email == "" {
        return  nil, fmt.Errorf("email is required")
    }
    if passhash == "" {
        return  nil, fmt.Errorf("passhash is required")
    }
    // TODO: add admin role
    user, err := s.userRepo.Create(&models.User{
        Username: username,
        Email: email,
        PassHash: passhash,
        Role: "user",
    })
    if err != nil {
        return nil, err
    }
    return user, nil
}
    
