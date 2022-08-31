package userstorage

import (
	"context"
	"errors"

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

func (u *userStorage) Insert(ctx context.Context, userName string, password string) (int, error) {
	sql := `INSERT INTO users(user_name,user_password)
			SELECT $1,$2
			WHERE NOT EXISTS (SELECT user_id FROM users WHERE user_name =$3) RETURNING user_id`
	var id int
	err := u.client.QueryRow(ctx, sql, userName, password, userName).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, errs.New(
				errs.Validation,
				errs.Code("user name already exists"),
				errs.Parameter("username"),
				err)
		}
		return -1, err
	}
	return id, nil
}

func (u *userStorage) Find(ctx context.Context, title string) (int, error) {
	sql := `SELECT vote_id FROM vote WHERE vote_title = $1`
	var id int
	err := u.client.QueryRow(ctx, sql, title).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, errs.New(
				errs.Validation,
				errs.Code("user name not found"),
				errs.Parameter("username"),
				err)
		}
		return -1, err
	}
	return id, nil
}
