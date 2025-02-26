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

func (s *RegisterService) Register(username, email, passhash string) error {
    if username == ""{
        return fmt.Errorf("username is required")
    }
    if email == "" {
        return fmt.Errorf("email is required")
    }
    if passhash == "" {
        return fmt.Errorf("passhash is required")
    }
    err := s.userRepo.Create(&models.User{Username: username, Email: email, PassHash: passhash})
    if err != nil {
        return err
    }
    return nil
}
    
