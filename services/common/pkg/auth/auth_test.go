package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestKeys() (string, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	privateBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	privatePem := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateBytes})
	privateBase64 := base64.StdEncoding.EncodeToString(privatePem)

	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	publicPem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicBytes})
	publicBase64 := base64.StdEncoding.EncodeToString(publicPem)

	return privateBase64, publicBase64, nil
}

func generateTestToken(privateKey *ecdsa.PrivateKey) (string, error) {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(5 * time.Minute).Unix(),
		"userID": 123,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(privateKey)
}

func TestNewVerifier(t *testing.T) {
	_, publicKeyBase64, err := generateTestKeys()
	assert.NoError(t, err)

	authService, err := NewVerifier(publicKeyBase64)
	assert.NoError(t, err)
	assert.NotNil(t, authService)
}

func TestVerifyJWTToken(t *testing.T) {
	privateBase64, publicBase64, err := generateTestKeys()
	assert.NoError(t, err)

	authService, err := NewVerifier(publicBase64)
	assert.NoError(t, err)

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateBase64)
	assert.NoError(t, err)

	privateBlock, _ := pem.Decode(privateKeyBytes)
	assert.NotNil(t, privateBlock)

	privateKey, err := x509.ParseECPrivateKey(privateBlock.Bytes)
	assert.NoError(t, err)

	token, err := generateTestToken(privateKey)
	assert.NoError(t, err)

	valid, err := authService.VerifyJWTToken(token)
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestVerifyExpiredJWTToken(t *testing.T) {
	privateBase64, publicBase64, err := generateTestKeys()
	assert.NoError(t, err)

	authService, err := NewVerifier(publicBase64)
	assert.NoError(t, err)

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateBase64)
	assert.NoError(t, err)

	privateBlock, _ := pem.Decode(privateKeyBytes)
	assert.NotNil(t, privateBlock)

	privateKey, err := x509.ParseECPrivateKey(privateBlock.Bytes)
	assert.NoError(t, err)

	claims := jwt.MapClaims{
		"exp":    time.Now().Add(-5 * time.Minute).Unix(), // Токен уже истёк
		"userID": 123,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := token.SignedString(privateKey)
	assert.NoError(t, err)

	valid, err := authService.VerifyJWTToken(signedToken)
	assert.Error(t, err)
	assert.False(t, valid)
}

