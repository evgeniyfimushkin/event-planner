package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var (
	ErrMissingToken = errors.New("access token is missing")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken = errors.New("invalid access token")
	ErrInvalidClaims = errors.New("cannot parse token claims")
	ErrInvalidExpClaim = errors.New("invalid exp claim in access token")
	ErrTokenExpired = errors.New("token is expired")
)

// Verifier - sevice, that contains publicKey and verifys jwt tokens
type Verifier struct {
	publicKey  *ecdsa.PublicKey
}

// NewVerifier takes a Base64-encoded EC256 public key,
// decodes them, and returns an instance of Verifier
func NewVerifier (publicKeyString string) (*Verifier, error) {

    publicKetBytes, err := base64.StdEncoding.DecodeString(publicKeyString)
    if err != nil {
        return nil, fmt.Errorf("failed to decode public key: %w", err)
    }

    publicBlock, _ := pem.Decode(publicKetBytes)
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

    return &Verifier{
        publicKey: publicKey,
    }, nil
}

// VerifyJWTToken takes accessToken as string and verified the signature
func (v *Verifier) VerifyJWTToken(accessToken string) (jwt.MapClaims, error) {
	if accessToken == "" {
		return nil, ErrMissingToken
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return v.publicKey, nil
	})
	if err != nil {
		// Если ошибка вызвана просроченностью токена, возвращаем ErrTokenExpired
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrInvalidExpClaim
	}
	if time.Now().Unix() > int64(exp) {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

