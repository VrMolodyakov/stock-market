package token

import (
	"fmt"
	"time"

	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/golang-jwt/jwt"
)

type tokenHandler struct {
	logger     *logging.Logger
	privateKey []byte
	publicKey  []byte
}

func NewTokenHandler(logger *logging.Logger) *tokenHandler {
	return &tokenHandler{logger: logger}
}

func (t *tokenHandler) CreateToken(ttl time.Duration, payload interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("couldn't parse private key: %w ", err)
	}
	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodES512, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("couldn't sign key due to %w", err)
	}
	return token, nil
}

func (t *tokenHandler) ValidateToken(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(t.publicKey)
	if err != nil {
		return "", fmt.Errorf("couldn't parse public key: %w ", err)
	}
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, signed := t.Method.(*jwt.SigningMethodRSA); !signed {
			return nil, fmt.Errorf("unexpected method - %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't parse : %w", err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token : %w", err)
	}
	return claims["sub"], nil
}
