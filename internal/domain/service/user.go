package service

import (
	"context"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
)

type UserStorage interface {
	Insert(ctx context.Context, username string, password string) (entity.User, error)
	Find(ctx context.Context, username string) (entity.User, error)
}

type userService struct {
	logger  *logging.Logger
	storage UserStorage
}

func NewUserStorage(logger *logging.Logger, storage UserStorage) *userService {
	return &userService{logger: logger, storage: storage}
}

func (u *userService) Create(ctx context.Context, username string, password string) (entity.User, error) {
	if username == "" {
		return entity.User{}, errs.New(errs.Validation, errs.Parameter("username"), errs.Code("empty username"))
	}
	if password == "" {
		return entity.User{}, errs.New(errs.Validation, errs.Parameter("password"), errs.Code("empty password"))
	}
	u.logger.Debug("create user with login = %v , password = %v", username, password)
	return u.storage.Insert(ctx, username, password)
}

func (u *userService) Get(ctx context.Context, username string) (entity.User, error) {
	if username == "" {
		return entity.User{}, errs.New(errs.Validation, errs.Parameter("username"), errs.Code("empty username"))
	}
	u.logger.Debugf("try to get user with username = %v", username)
	return u.storage.Find(ctx, username)
}
