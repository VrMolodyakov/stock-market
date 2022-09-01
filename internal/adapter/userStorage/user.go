package userstorage

import (
	"context"
	"errors"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type DbClient interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type userStorage struct {
	logger logging.Logger
	client DbClient
}

func New(logger logging.Logger, client DbClient) *userStorage {
	return &userStorage{logger: logger, client: client}
}

func (u *userStorage) Insert(ctx context.Context, username string, password string) (entity.User, error) {
	sql := `INSERT INTO users(u_name,u_password)
			SELECT $1,$2
			WHERE NOT EXISTS (SELECT u_id FROM users WHERE u_name =$3) RETURNING u_id,u_name,u_password`
	var user entity.User
	err := u.client.QueryRow(ctx, sql, username, password, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errs.New(
				errs.Validation,
				errs.Code("user name already exists"),
				errs.Parameter("username"),
				err)
		}
		return entity.User{}, err
	}
	return user, nil
}

func (u *userStorage) Find(ctx context.Context, username string) (entity.User, error) {
	sql := `SELECT u_id FROM users WHERE u_name = $1`
	var user entity.User
	err := u.client.QueryRow(ctx, sql, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errs.New(
				errs.Validation,
				errs.Code("user name not found"),
				errs.Parameter("username"),
				err)
		}
		return entity.User{}, err
	}
	return user, nil
}
