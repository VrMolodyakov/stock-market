package service

import (
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
)

type TokenStorage interface {
	Set(refreshToken string, userId int, expireAt time.Duration) error
	Get(refreshToken string) (int, error)
	Delete(refreshToken string) error
}

type tokenService struct {
	logger  *logging.Logger
	storage TokenStorage
}

func NewTokenService(storage TokenStorage, logger *logging.Logger) *tokenService {
	return &tokenService{storage: storage, logger: logger}
}

func (t *tokenService) Save(refreshToken string, userId int, expireAt time.Duration) error {
	if len(refreshToken) == 0 {
		return errs.New(errs.Validation, errs.Code("refresh token is empty"), errs.Parameter("refresh token"))
	}
	if userId < 0 {
		return errs.New(errs.Validation, errs.Code("user id can't be less than zero"), errs.Parameter("user id"))
	}
	return t.storage.Set(refreshToken, userId, expireAt)
}

func (t *tokenService) Find(refreshToken string) (int, error) {
	if len(refreshToken) == 0 {
		return -1, errs.New(errs.Validation, errs.Code("refresh token is empty"), errs.Parameter("refresh token"))
	}
	return t.storage.Get(refreshToken)
}

func (t *tokenService) Remove(refreshToken string) error {
	if len(refreshToken) == 0 {
		return errs.New(errs.Validation, errs.Code("refresh token is empty"), errs.Parameter("refresh token"))
	}
	return t.storage.Delete(refreshToken)
}
