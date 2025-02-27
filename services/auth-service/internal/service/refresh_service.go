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

type RefreshService struct {
    userRepo *repository.UserRepository
    privateKey *ecdsa.PrivateKey
    publicKey *ecdsa.PublicKey
}

func NewRefreshService (userRepo *repository.UserRepository, privateKeyString string, publicKeyString string) (*RefreshService, error) {
    privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
    if err != nil {
        return nil, fmt.Errorf("failed to decode private key: %w", err)
	}
    
    privateBlock, _ := pem.Decode(privateKeyBytes)
    if privateBlock == nil {
        return nil, errors.New("failed to parse PEM block containing the private key")
    }
    privateKey, err := x509.ParseECPrivateKey(privateBlock.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parce EC private key: %w", err)
    }

    publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyString)
    if err != nil {
        return nil, fmt.Errorf("failed to decode public key: %w", err)
    }
    
    publicBlock, _ := pem.Decode(publicKeyBytes)
    if publicBlock == nil {
        return nil, errors.New("failed to parse PEM block containing the public key")
    }
    
    publicKeyInterface, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse EC public key: %w", err)
    }

    publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
    if !ok {
        return nil, errors.New("public key is not a valid ECDSA key")
    } 

    return &RefreshService{
        userRepo: userRepo,
        privateKey: privateKey,
        publicKey:  publicKey,
    }, nil
}

func (r *RefreshService) Refresh(refreshToken string) (string, error) {
    token, err := jwt.Parse(refreshToken, func (token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return r.publicKey, nil
    })
    if err != nil {
        return "", fmt.Errorf("invalid refresh token: %w", err)
    }

    if !token.Valid {
        return "", errors.New("invalid refresh token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", errors.New("cannot parse token claims")
    }
    exp, ok := claims["exp"].(float64)
    if !ok {
        return "", errors.New("invalid exp claim in refresh token")
    }

    if time.Now().Unix() > int64(exp) {
        return "", errors.New("refresh token is too old")
    }

    userIDFloat, ok := claims["userID"].(float64)
    if !ok {
        return "", errors.New("invalid userid claim in refresh token")
    }
    
    userID := int(userIDFloat)


    user, err := r.userRepo.GetByID(userID)
    if err != nil {
        return "", fmt.Errorf("authentication failed: %w", err)
    }

    return r.generateAccessJWT(user, 15 * time.Minute)
}

func (r *RefreshService) generateAccessJWT(user *models.User, tokenTTL time.Duration) (string, error) {
    token := jwt.New(jwt.SigningMethodES256)
    claims := token.Claims.(jwt.MapClaims)
    claims["userID"] = user.ID
    claims["username"] = user.Username
    claims["email"] = user.Email
    claims["role"] = user.Role
    claims["exp"] = time.Now().Add(tokenTTL).Unix()
    tokenString, err := token.SignedString(r.privateKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
