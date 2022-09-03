package token

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/golang-jwt/jwt"
)

type TokenPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

type tokenHandler struct {
	logger      *logging.Logger
	accessPair  TokenPair
	refreshPair TokenPair
}

func NewTokenHandler(logger *logging.Logger, accessPair TokenPair, refreshPair TokenPair) *tokenHandler {
	return &tokenHandler{logger: logger, accessPair: accessPair, refreshPair: refreshPair}
}

func (t *tokenHandler) CreateAccessToken(ttl time.Duration, payload interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.accessPair.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("couldn't parse private key: %w ", err)
	}
	return create(ttl, payload, key)
}

func (t *tokenHandler) CreateRefreshToken(ttl time.Duration, payload interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.refreshPair.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("couldn't parse private key: %w ", err)
	}
	return create(ttl, payload, key)
}

func create(ttl time.Duration, payload interface{}, key *rsa.PrivateKey) (string, error) {
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

func (t *tokenHandler) ValidateAccessToken(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(t.accessPair.PublicKey)
	if err != nil {
		return "", fmt.Errorf("couldn't parse public key: %w ", err)
	}
	return validate(token, key)
}

func (t *tokenHandler) ValidateRefreshToken(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(t.refreshPair.PublicKey)
	if err != nil {
		return "", fmt.Errorf("couldn't parse public key: %w ", err)
	}
	return validate(token, key)

}

func validate(token string, key *rsa.PublicKey) (interface{}, error) {
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
