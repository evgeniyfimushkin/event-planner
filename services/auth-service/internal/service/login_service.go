package service

import (
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginService struct {
    userRepo *repository.UserRepository
    privateKey *ecdsa.PrivateKey
    tokenTTL time.Duration
}

func NewLoginService(userRepo *repository.UserRepository, secret string, tokenTTL time.Duration) (*LoginService, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %w", err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	return &LoginService{
		userRepo:  userRepo,
		privateKey: privateKey,
        tokenTTL: tokenTTL,
	}, nil
}

func (s *LoginService) Login(username, passhash string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	if user.PassHash != passhash {
		return "", fmt.Errorf("incorrect password: %w", err)
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *LoginService) generateJWT(user *models.User)(string, error) {
    token := jwt.New(jwt.SigningMethodES256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = user.Username
    claims["email"] = user.Email
    claims["exp"] = time.Now().Add(s.tokenTTL).Unix()
    claims["iat"] = time.Now().Unix()
    tokenString, err := token.SignedString(s.privateKey)
    if err != nil {
        return "", err
    }
    return tokenString, err
}
