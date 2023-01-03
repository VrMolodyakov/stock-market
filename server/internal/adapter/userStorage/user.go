package userstorage

import (
	"context"
	"errors"
	"time"

	"github.com/VrMolodyakov/stock-market/internal/domain/entity"
	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type DbClient interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type userStorage struct {
	logger *logging.Logger
	client DbClient
}

func New(logger *logging.Logger, client DbClient) *userStorage {
	return &userStorage{logger: logger, client: client}
}

func (u *userStorage) Insert(ctx context.Context, username string, password string) (entity.User, error) {
	sql := `INSERT INTO users(u_name,u_password,create_at)
			SELECT $1,$2,$3
			WHERE NOT EXISTS (SELECT u_id FROM users WHERE u_name =$4) RETURNING u_id,u_name,u_password,create_at`
	var user entity.User
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)
	err := u.client.QueryRow(ctx, sql, username, password, dt, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreateAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errs.New(
				errs.Validation,
				errs.Code("user already exists"),
				errs.Parameter("username"),
				err)
		}
		return entity.User{}, err
	}
	return user, nil
}

func (u *userStorage) Find(ctx context.Context, username string) (entity.User, error) {
	sql := `SELECT u_id,u_name,u_password,create_at FROM users WHERE u_name = $1`
	var user entity.User
	err := u.client.QueryRow(ctx, sql, username).Scan(&user.Id, &user.Username, &user.Password, &user.CreateAt)
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

func (u *userStorage) FindById(ctx context.Context, id int) (entity.User, error) {
	sql := `SELECT u_id,u_name,u_password,create_at FROM users WHERE u_id = $1`
	var user entity.User
	err := u.client.QueryRow(ctx, sql, id).Scan(&user.Id, &user.Username, &user.Password, &user.CreateAt)
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
