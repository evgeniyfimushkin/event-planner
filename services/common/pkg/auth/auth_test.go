package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateECDSAKeys генерирует пару ключей и возвращает приватный, публичный ключ и строку с Base64-encoded PEM представлением публичного ключа.
func generateECDSAKeys(t *testing.T) (*ecdsa.PrivateKey, *ecdsa.PublicKey, string) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ECDSA key: %v", err)
	}
	pub := &priv.PublicKey

	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	encoded := base64.StdEncoding.EncodeToString(pemBlock)
	return priv, pub, encoded
}

// createToken создаёт JWT токен с заданными claims, используя указанный метод подписи.
func createToken(t *testing.T, priv *ecdsa.PrivateKey, claims jwt.MapClaims, method jwt.SigningMethod) string {
	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return tokenString
}

// TestNewVerifier_Valid проверяет успешное создание Verifier с корректным публичным ключом.
func TestNewVerifier_Valid(t *testing.T) {
	_, _, pubKeyStr := generateECDSAKeys(t)
	verif, err := NewVerifier(pubKeyStr)
	if err != nil {
		t.Fatalf("expected valid verifier, got error: %v", err)
	}
	if verif == nil {
		t.Fatal("expected non-nil verifier")
	}
}

// TestNewVerifier_Invalid проверяет создание Verifier с некорректной строкой публичного ключа.
func TestNewVerifier_Invalid(t *testing.T) {
	_, err := NewVerifier("invalid-base64")
	if err == nil {
		t.Fatal("expected error for invalid base64 input")
	}
}

// TestVerifyJWTToken_EmptyToken проверяет, что при отсутствии токена возвращается ErrMissingToken.
func TestVerifyJWTToken_EmptyToken(t *testing.T) {
	_, _, pubKeyStr := generateECDSAKeys(t)
	verif, err := NewVerifier(pubKeyStr)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}
	_, err = verif.VerifyJWTToken("")
	if !errors.Is(err, ErrMissingToken) {
		t.Fatalf("expected ErrMissingToken, got %v", err)
	}
}

// TestVerifyJWTToken_InvalidSigningMethod проверяет обработку токена с неподходящим методом подписи.
func TestVerifyJWTToken_InvalidSigningMethod(t *testing.T) {
	_, _, pubKeyStr := generateECDSAKeys(t)
	verif, err := NewVerifier(pubKeyStr)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}
	// Создаём токен с неподходящим методом подписи (HS256 вместо ECDSA)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	_, err = verif.VerifyJWTToken(tokenString)
	if err == nil || !errors.Is(err, ErrUnexpectedSigningMethod) {
		t.Fatalf("expected ErrUnexpectedSigningMethod, got %v", err)
	}
}

// TestVerifyJWTToken_ExpiredToken проверяет, что при истёкшем сроке действия токена возвращается ErrTokenExpired.
func TestVerifyJWTToken_ExpiredToken(t *testing.T) {
	priv, _, pubKeyStr := generateECDSAKeys(t)
	verif, err := NewVerifier(pubKeyStr)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(-1 * time.Hour).Unix(), // токен просрочен
	}
	tokenString := createToken(t, priv, claims, jwt.SigningMethodES256)
	_, err = verif.VerifyJWTToken(tokenString)
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

// TestVerifyJWTToken_Valid проверяет успешную валидацию корректного токена и извлечение claims.
func TestVerifyJWTToken_Valid(t *testing.T) {
	priv, _, pubKeyStr := generateECDSAKeys(t)
	verif, err := NewVerifier(pubKeyStr)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	expectedValue := "test-value"
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
		"custom": expectedValue,
	}
	tokenString := createToken(t, priv, claims, jwt.SigningMethodES256)
	retClaims, err := verif.VerifyJWTToken(tokenString)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if retClaims["custom"] != expectedValue {
		t.Fatalf("expected custom claim %v, got %v", expectedValue, retClaims["custom"])
	}
}

// TestVerifyJWTToken_InvalidTokenFormat проверяет обработку токена с некорректной структурой.
func TestVerifyJWTToken_InvalidTokenFormat(t *testing.T) {
	_, _, pubKeyStr := generateECDSAKeys(t)
	verif, err := NewVerifier(pubKeyStr)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}

	_, err = verif.VerifyJWTToken("not-a-valid-token")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
}

