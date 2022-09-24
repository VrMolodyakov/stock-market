package service

import (
	"context"

	"github.com/VrMolodyakov/stock-market/internal/domain/entity"
	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
)

type UserStorage interface {
	Insert(ctx context.Context, username string, password string) (entity.User, error)
	Find(ctx context.Context, username string) (entity.User, error)
	FindById(ctx context.Context, id int) (entity.User, error)
}

type userService struct {
	logger  *logging.Logger
	storage UserStorage
}

func NewUserService(logger *logging.Logger, storage UserStorage) *userService {
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

func (u *userService) GetById(ctx context.Context, id int) (entity.User, error) {
	if id < 0 {
		return entity.User{}, errs.New(errs.Validation, errs.Parameter("id"), errs.Code("id less than zero"))
	}
	u.logger.Debugf("try to get user with id = %v", id)
	return u.storage.FindById(ctx, id)
}
