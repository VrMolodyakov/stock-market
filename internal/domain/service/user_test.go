package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/internal/domain/service/mocks"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreare(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserStorage(ctrl)
	defer ctrl.Finish()
	type mock func(username string, password string) *userService
	type args struct {
		username string
		password string
	}

	testCases := []struct {
		title    string
		mockCall mock
		input    args
		isError  bool
		want     entity.User
	}{
		{
			title: "Success create and return nil",
			mockCall: func(username string, password string) *userService {
				logger := logging.GetLogger("debug")
				user := entity.User{Username: username, Id: 1, Password: password}
				userRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
				return NewUserService(logger, userRepo)
			},
			input:   args{username: "username", password: "password"},
			want:    entity.User{Username: "username", Password: "password"},
			isError: false,
		},
		{
			title: "Empty username and return error",
			mockCall: func(username string, password string) *userService {
				logger := logging.GetLogger("debug")
				return NewUserService(logger, userRepo)
			},
			input:   args{username: "", password: "password"},
			isError: true,
		},
		{
			title: "Empty password and return error",
			mockCall: func(username string, password string) *userService {
				logger := logging.GetLogger("debug")
				return NewUserService(logger, userRepo)
			},
			input:   args{username: "username", password: ""},
			isError: true,
		},
		{
			title: "Internal db error",
			mockCall: func(username string, password string) *userService {
				logger := logging.GetLogger("debug")
				userRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.User{}, errors.New("internal db error"))
				return NewUserService(logger, userRepo)
			},
			input:   args{username: "username", password: "password"},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			userService := test.mockCall(test.input.username, test.input.password)
			got, err := userService.Create(context.Background(), test.input.username, test.input.password)
			if !test.isError {
				assert.Equal(t, test.want.Username, got.Username)
				assert.Equal(t, test.want.Password, got.Password)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserStorage(ctrl)
	defer ctrl.Finish()
	type mock func(username string) *userService
	type args struct {
		username string
	}

	testCases := []struct {
		title    string
		mockCall mock
		input    args
		isError  bool
		want     entity.User
	}{
		{
			title: "Success find and return nil",
			mockCall: func(username string) *userService {
				logger := logging.GetLogger("debug")
				user := entity.User{Username: username, Id: 1, Password: "password", CreateAt: time.Now()}
				userRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(user, nil)
				return NewUserService(logger, userRepo)
			},
			input:   args{username: "username"},
			want:    entity.User{Username: "username"},
			isError: false,
		},
		{
			title: "Empty username and return error",
			mockCall: func(username string) *userService {
				logger := logging.GetLogger("debug")
				return NewUserService(logger, userRepo)
			},
			input:   args{username: ""},
			isError: true,
		},
		{
			title: "Internal db error",
			mockCall: func(username string) *userService {
				logger := logging.GetLogger("debug")
				userRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(entity.User{}, errors.New("internal db error"))
				return NewUserService(logger, userRepo)
			},
			input:   args{username: "username"},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			userService := test.mockCall(test.input.username)
			got, err := userService.Get(context.Background(), test.input.username)
			if !test.isError {
				assert.Equal(t, test.want.Username, got.Username)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestGetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserStorage(ctrl)
	defer ctrl.Finish()
	type mock func(userId int) *userService
	type args struct {
		userId int
	}

	testCases := []struct {
		title    string
		mockCall mock
		input    args
		isError  bool
		want     entity.User
	}{
		{
			title: "Success find and return nil",
			mockCall: func(userId int) *userService {
				logger := logging.GetLogger("debug")
				user := entity.User{Username: "username", Id: userId, Password: "password", CreateAt: time.Now()}
				userRepo.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(user, nil)
				return NewUserService(logger, userRepo)
			},
			input:   args{userId: 1},
			want:    entity.User{Username: "username"},
			isError: false,
		},
		{
			title: "Id less than zero and return error",
			mockCall: func(userId int) *userService {
				logger := logging.GetLogger("debug")
				return NewUserService(logger, userRepo)
			},
			input:   args{userId: -1},
			isError: true,
		},
		{
			title: "Id less than zero and return error",
			mockCall: func(userId int) *userService {
				logger := logging.GetLogger("debug")
				userRepo.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(entity.User{}, errors.New("internal db error"))
				return NewUserService(logger, userRepo)
			},
			input:   args{userId: 1},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			userService := test.mockCall(test.input.userId)
			got, err := userService.GetById(context.Background(), test.input.userId)
			if !test.isError {
				assert.Equal(t, test.want.Username, got.Username)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
