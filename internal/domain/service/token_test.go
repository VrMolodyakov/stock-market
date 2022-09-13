package service

import (
	"errors"
	"testing"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/service/mocks"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	tokenRepo := mocks.NewMockTokenStorage(ctrl)
	defer ctrl.Finish()
	type mock func() *tokenService
	type args struct {
		refreshToken string
		userId       int
		expireAt     time.Duration
	}

	testCases := []struct {
		title    string
		mockCall mock
		input    args
		isError  bool
	}{
		{
			title: "Success save token",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				tokenRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "refresh token", userId: 1, expireAt: 5 * time.Second},
			isError: false,
		},
		{
			title: "Refresh token is empty and return error",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "", userId: 1, expireAt: 5 * time.Second},
			isError: true,
		},
		{
			title: "User id is less than zero and return error",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "refresh token", userId: -1, expireAt: 5 * time.Second},
			isError: true,
		},
		{
			title: "Internal db error error",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				tokenRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("internal db error"))
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "refresh token", userId: 1, expireAt: 5 * time.Second},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			tokenService := test.mockCall()
			err := tokenService.Save(test.input.refreshToken, test.input.userId, test.input.expireAt)
			if !test.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestFind(t *testing.T) {
	ctrl := gomock.NewController(t)
	tokenRepo := mocks.NewMockTokenStorage(ctrl)
	defer ctrl.Finish()
	type mock func() *tokenService
	type args struct {
		refreshToken string
	}

	testCases := []struct {
		title    string
		mockCall mock
		input    args
		isError  bool
		want     int
	}{
		{
			title: "Success save token",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				tokenRepo.EXPECT().Get(gomock.Any()).Return(1, nil)
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "refresh token"},
			want:    1,
			isError: false,
		},
		{
			title: "Refresh token is empty and return error",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: ""},
			want:    -1,
			isError: true,
		},
		{
			title: "Internal db error error",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				tokenRepo.EXPECT().Get(gomock.Any()).Return(-1, errors.New("internal db error"))
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "refresh token"},
			want:    -1,
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			tokenService := test.mockCall()
			got, err := tokenService.Find(test.input.refreshToken)
			if !test.isError {
				assert.Equal(t, test.want, got)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	tokenRepo := mocks.NewMockTokenStorage(ctrl)
	defer ctrl.Finish()
	type mock func() *tokenService
	type args struct {
		refreshToken string
	}

	testCases := []struct {
		title    string
		mockCall mock
		input    args
		isError  bool
	}{
		{
			title: "Success save token",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				tokenRepo.EXPECT().Delete(gomock.Any()).Return(nil)
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "refresh token"},
			isError: false,
		},
		{
			title: "Refresh token is empty and return error",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: ""},
			isError: true,
		},
		{
			title: "Internal db error error",
			mockCall: func() *tokenService {
				logger := logging.GetLogger("debug")
				tokenRepo.EXPECT().Delete(gomock.Any()).Return(errors.New("internal db error"))
				return NewTokenService(tokenRepo, logger)
			},
			input:   args{refreshToken: "refresh token"},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			tokenService := test.mockCall()
			err := tokenService.Remove(test.input.refreshToken)
			if !test.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
