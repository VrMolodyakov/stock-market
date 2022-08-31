package service

import (
	"context"

	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
)

type UserStorage interface {
	Insert(ctx context.Context, userName string, password string) (int, error)
	Find(ctx context.Context, title string) (int, error)
}

type userService struct {
	logger  *logging.Logger
	storage UserStorage
}

func NewUserStorage(logger *logging.Logger, storage UserStorage) *userService {
	return &userService{logger: logger, storage: storage}
}

func (u *userService) Create()
