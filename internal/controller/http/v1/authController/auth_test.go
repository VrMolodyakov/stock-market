package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/controller/http/v1/authController/mocks"
	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignUpUser(t *testing.T) {
	cntr := gomock.NewController(t)
	mockUserService := mocks.NewMockUserService(cntr)
	mockTokenHandler := mocks.NewMockTokenHandler(cntr)
	mockTokenService := mocks.NewMockTokenService(cntr)
	now := time.Now()
	inputTime := now.Format(time.RFC3339)
	authController := NewAuthController(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, 15, 15)
	type mockCall func()
	testCases := []struct {
		title        string
		mock         mockCall
		inputRequest string
		want         string
		expectedCode int
	}{
		{
			title: "sign up and 201 response",
			mock: func() {
				user := entity.User{Username: "username", Password: "hashedpassword", CreateAt: now, Id: 1}
				mockUserService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
			},
			inputRequest: `{"username":"username","password":"password"}`,
			want:         fmt.Sprintf("{\"data\":{\"user\":{\"username\":\"username\",\"password\":\"hashedpassword\",\"create_at\":\"%v\"}},\"status\":\"success\"}", inputTime),
			expectedCode: 201,
		},
		{
			title: "wrong input request and 500 response",
			mock: func() {
			},
			inputRequest: `wrong data`,
			want:         "",
			expectedCode: 500,
		},
		{
			title: "Cannot create user and 500 response",
			mock: func() {
				mockUserService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.User{}, errs.New(errs.Database))
			},
			inputRequest: `{"username":"username","password":"password"}`,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			expectedCode: 500,
		},
		{
			title: "Internal service error and 500 response",
			mock: func() {
				mockUserService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.User{}, errs.New(errs.Internal))
			},
			inputRequest: `{"username":"username","password":"password"}`,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			expectedCode: 500,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			router := gin.Default()
			router.POST("/register", authController.SignUpUser)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(test.inputRequest))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.want, recorder.Body.String())
			assert.Equal(t, test.expectedCode, recorder.Code)
		})
	}
}

func TestSignInUser(t *testing.T) {
	cntr := gomock.NewController(t)
	mockUserService := mocks.NewMockUserService(cntr)
	mockTokenHandler := mocks.NewMockTokenHandler(cntr)
	mockTokenService := mocks.NewMockTokenService(cntr)
	authController := NewAuthController(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, 15, 15)
	type mockCall func(accessToken string, refreshToken string)
	testCases := []struct {
		title        string
		mock         mockCall
		inputRequest string
		want         string
		expectedCode int
		tokens       []string
		isError      bool
	}{
		{
			title: "sign up and 201 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, nil)
				mockTokenHandler.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(refreshToken, nil)
				mockTokenService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 200,
			want:         "{\"access_token\":\"encodedAccessToken\",\"status\":\"success\"}",
			isError:      false,
			tokens:       []string{"encodedAccessToken", "encodedRefreshToken"},
		},
		{
			title: "user not found",
			mock: func(accessToken string, refreshToken string) {
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entity.User{}, errs.New(errs.Validation, errs.Parameter("username")))
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 400,
			want:         "\"username\"",
			tokens:       []string{"", ""},
			isError:      true,
		},
		{
			title: "cannot create access token and 500 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, errors.New("internal token service error"))
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 500,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			tokens:       []string{"", ""},
			isError:      true,
		},
		{
			title: "cannot create refresh token and 500 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, nil)
				mockTokenHandler.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(refreshToken, errors.New("internal token service error"))
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 500,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			tokens:       []string{"encodedAccessToken", ""},
			isError:      true,
		},
		{
			title: "cannot save refresh token and 500 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, nil)
				mockTokenHandler.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(refreshToken, nil)
				mockTokenService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errs.New(errs.Internal))

			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 500,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			isError:      true,
			tokens:       []string{"encodedAccessToken", "encodedRefreshToken"},
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock(test.tokens[0], test.tokens[1])
			router := gin.Default()
			router.POST("/login", authController.SignInUser)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(test.inputRequest))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			if !test.isError {
				coockies := recorder.Result().Cookies()
				for _, c := range coockies {
					if c.Name == "access_token" {
						assert.Equal(t, test.tokens[0], c.Value)
					} else if c.Name == "refresh_token" {
						assert.Equal(t, test.tokens[1], c.Value)
					}
				}

			}
			assert.Equal(t, test.want, recorder.Body.String())
			assert.Equal(t, test.expectedCode, recorder.Code)
		})
	}
}

func TestRefreshAccessToken(t *testing.T) {
	cntr := gomock.NewController(t)
	mockUserService := mocks.NewMockUserService(cntr)
	mockTokenHandler := mocks.NewMockTokenHandler(cntr)
	mockTokenService := mocks.NewMockTokenService(cntr)
	authController := NewAuthController(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, 15, 15)
	type mockCall func(accessToken string, refreshToken string)
	testCases := []struct {
		title        string
		mock         mockCall
		inputRequest string
		want         string
		expectedCode int
		tokens       []string
		isError      bool
	}{
		{
			title: "sign up and 201 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, nil)
				mockTokenHandler.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(refreshToken, nil)
				mockTokenService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 200,
			want:         "{\"access_token\":\"encodedAccessToken\",\"status\":\"success\"}",
			isError:      false,
			tokens:       []string{"encodedAccessToken", "encodedRefreshToken"},
		},
		{
			title: "user not found",
			mock: func(accessToken string, refreshToken string) {
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entity.User{}, errs.New(errs.Validation, errs.Parameter("username")))
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 400,
			want:         "\"username\"",
			tokens:       []string{"", ""},
			isError:      true,
		},
		{
			title: "cannot create access token and 500 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, errors.New("internal token service error"))
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 500,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			tokens:       []string{"", ""},
			isError:      true,
		},
		{
			title: "cannot create refresh token and 500 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, nil)
				mockTokenHandler.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(refreshToken, errors.New("internal token service error"))
			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 500,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			tokens:       []string{"encodedAccessToken", ""},
			isError:      true,
		},
		{
			title: "cannot save refresh token and 500 response",
			mock: func(accessToken string, refreshToken string) {
				user := entity.User{Username: "username", Password: "$2a$10$EY89/z9fLxDtT0V18CMYje2K5.q28PPkbaQAuvLJ8pJJF.nElg.r6", CreateAt: time.Now(), Id: 1}
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(user, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, nil)
				mockTokenHandler.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(refreshToken, nil)
				mockTokenService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errs.New(errs.Internal))

			},
			inputRequest: `{"username":"username","password":"my_password"}`,
			expectedCode: 500,
			want:         "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			isError:      true,
			tokens:       []string{"encodedAccessToken", "encodedRefreshToken"},
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock(test.tokens[0], test.tokens[1])
			router := gin.Default()
			router.POST("/login", authController.SignInUser)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(test.inputRequest))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			if !test.isError {
				coockies := recorder.Result().Cookies()
				for _, c := range coockies {
					if c.Name == "access_token" {
						assert.Equal(t, test.tokens[0], c.Value)
					} else if c.Name == "refresh_token" {
						assert.Equal(t, test.tokens[1], c.Value)
					}
				}

			}
			assert.Equal(t, test.want, recorder.Body.String())
			assert.Equal(t, test.expectedCode, recorder.Code)
		})
	}
}
