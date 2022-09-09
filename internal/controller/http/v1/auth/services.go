package v1

import (
	"context"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
)

type UserService interface {
	Create(ctx context.Context, username string, password string) (entity.User, error)
	Get(ctx context.Context, username string) (entity.User, error)
}

type TokenHandler interface {
	CreateAccessToken(ttl time.Duration, payload interface{}) (string, error)
	CreateRefreshToken(ttl time.Duration, payload interface{}) (string, error)
	ValidateRefreshToken(token string) error
}

type TokenService interface {
	Save(refreshToken string, userId int, expireAt time.Duration) error
	Find(refreshToken string) (int, error)
	Remove(refreshToken string) error
}
