package userstorage

import (
	"context"
	"testing"
	"time"

	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
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
	testCases := []struct {
		title string
		mock  mockCall
		args  args
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
					gomock.Any()).Return(row, nil)
			},
			args: args{username: "username", password: "password"},
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			mockClient.Insert(context.Background(), test.args.username, test.args.password)
		})
	}
}
