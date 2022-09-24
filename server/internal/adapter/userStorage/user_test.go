package userstorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

type userEntityRow struct {
	Id       int
	Username string
	Password string
	CreateAt time.Time
	Err      error
}

func (this userEntityRow) Scan(dest ...interface{}) error {
	if this.Err != nil {
		return this.Err
	}
	id := dest[0].(*int)
	username := dest[1].(*string)
	password := dest[2].(*string)
	createAt := dest[3].(*time.Time)
	*id = this.Id
	*username = this.Username
	*password = this.Password
	*createAt = this.CreateAt
	return nil
}

func TestInsertUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	mockClient := userStorage{client: mockPool, logger: logger}
	type mockCall func()
	type args struct {
		username string
		password string
	}
	type want struct {
		id       int
		username string
	}
	testCases := []struct {
		title   string
		mock    mockCall
		args    args
		want    want
		isError bool
	}{
		{
			title: "Should Insert new user",
			mock: func() {
				row := userEntityRow{Id: 1, Username: "username", Password: "password", CreateAt: time.Now()}
				mockPool.EXPECT().QueryRow(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(row)
			},
			args:    args{username: "username", password: "password"},
			want:    want{id: 1, username: "username"},
			isError: false,
		},
		{
			title: "Internal error",
			mock: func() {
				row := userEntityRow{Err: errors.New("internal error")}
				mockPool.EXPECT().QueryRow(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(row)
			},
			args:    args{username: "username", password: "password"},
			isError: true,
		},
		{
			title: "No row",
			mock: func() {
				row := userEntityRow{Err: pgx.ErrNoRows}
				mockPool.EXPECT().QueryRow(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(row)
			},
			args:    args{username: "username", password: "password"},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := mockClient.Insert(context.Background(), test.args.username, test.args.password)
			if !test.isError {
				assert.Equal(t, test.want.id, got.Id)
				assert.Equal(t, test.want.username, got.Username)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestFindUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	mockClient := userStorage{client: mockPool, logger: logger}
	type mockCall func()
	type want struct {
		id       int
		username string
	}
	testCases := []struct {
		title   string
		mock    mockCall
		args    string
		want    want
		isError bool
	}{
		{
			title: "Should Insert new user",
			mock: func() {
				row := userEntityRow{Id: 1, Username: "username", Password: "password", CreateAt: time.Now()}
				mockPool.EXPECT().QueryRow(
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(row)
			},
			args:    "username",
			want:    want{id: 1, username: "username"},
			isError: false,
		},
		{
			title: "Internal error",
			mock: func() {
				row := userEntityRow{Err: pgx.ErrNoRows}
				mockPool.EXPECT().QueryRow(
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(row)
			},
			args:    "username",
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := mockClient.Find(context.Background(), test.args)
			if !test.isError {
				assert.Equal(t, test.want.id, got.Id)
				assert.Equal(t, test.want.username, got.Username)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestFindByIdUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	mockClient := userStorage{client: mockPool, logger: logger}
	type mockCall func()
	type want struct {
		id       int
		username string
	}
	testCases := []struct {
		title   string
		mock    mockCall
		args    int
		want    want
		isError bool
	}{
		{
			title: "Should Insert new user",
			mock: func() {
				row := userEntityRow{Id: 1, Username: "username", Password: "password", CreateAt: time.Now()}
				mockPool.EXPECT().QueryRow(
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(row)
			},
			args:    1,
			want:    want{id: 1, username: "username"},
			isError: false,
		},
		{
			title: "Internal error",
			mock: func() {
				row := userEntityRow{Err: pgx.ErrNoRows}
				mockPool.EXPECT().QueryRow(
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(row)
			},
			args:    1,
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := mockClient.FindById(context.Background(), test.args)
			if !test.isError {
				assert.Equal(t, test.want.id, got.Id)
				assert.Equal(t, test.want.username, got.Username)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
