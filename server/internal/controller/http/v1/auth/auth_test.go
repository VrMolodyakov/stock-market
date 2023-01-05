package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VrMolodyakov/stock-market/internal/controller/http/v1/auth/mocks"
	"github.com/VrMolodyakov/stock-market/internal/domain/entity"
	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
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
	authHandler := NewAuthHandler(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, "localhost", 15, 15)
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
			want:         "\"incorrect data format\"",
			expectedCode: 400,
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
			router.POST("/register", authHandler.SignUpUser)
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
	authHandler := NewAuthHandler(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, "localhost", 15, 15)
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
				mockUserService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entity.User{}, errs.New(errs.Validation, errs.Code("username")))
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
			router.POST("/login", authHandler.SignInUser)
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
	authHandler := NewAuthHandler(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, "localhost", 15, 15)
	type mockCall func(recorder *httptest.ResponseRecorder, userId int, accessToken string)
	type args struct {
		acessToken string
		userId     int
	}
	testCases := []struct {
		title          string
		mock           mockCall
		args           args
		wantedResponse string
		expectedCode   int
		wantedTokens   []string
		isError        bool
	}{
		{
			title: "sign up and 200 response",
			args:  args{acessToken: "newAccessToken", userId: 1},
			mock: func(recorder *httptest.ResponseRecorder, userId int, accessToken string) {
				http.SetCookie(recorder, &http.Cookie{Name: "refresh_token", Value: "encodedRefreshToken"})
				mockTokenHandler.EXPECT().ValidateRefreshToken(gomock.Any()).Return(nil)
				mockTokenService.EXPECT().Find(gomock.Any()).Return(userId, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return(accessToken, nil)

			},
			expectedCode:   200,
			wantedTokens:   []string{"newAccessToken", "encodedRefreshToken"},
			wantedResponse: "{\"access_token\":\"newAccessToken\",\"status\":\"success\"}",
			isError:        false,
		},
		{
			title: "cannot validate token and 403 response",
			args:  args{acessToken: "", userId: 0},
			mock: func(recorder *httptest.ResponseRecorder, userId int, accessToken string) {
				http.SetCookie(recorder, &http.Cookie{Name: "refresh_token", Value: "encodedRefreshToken"})
				mockTokenHandler.EXPECT().ValidateRefreshToken(gomock.Any()).Return(errors.New("validation error"))

			},
			expectedCode:   403,
			wantedTokens:   []string{},
			wantedResponse: "\"validation error\"",
			isError:        true,
		},
		{
			title: "cannot find refresh token and 403 response",
			args:  args{acessToken: "", userId: 0},
			mock: func(recorder *httptest.ResponseRecorder, userId int, accessToken string) {
				http.SetCookie(recorder, &http.Cookie{Name: "refresh_token", Value: "encodedRefreshToken"})
				mockTokenHandler.EXPECT().ValidateRefreshToken(gomock.Any()).Return(nil)
				mockTokenService.EXPECT().Find(gomock.Any()).Return(-1, errs.New(errs.Database))

			},
			expectedCode:   500,
			wantedTokens:   []string{},
			wantedResponse: "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			isError:        true,
		},
		{
			title: "cannot create new access token and 500 response",
			args:  args{acessToken: "", userId: 0},
			mock: func(recorder *httptest.ResponseRecorder, userId int, accessToken string) {
				http.SetCookie(recorder, &http.Cookie{Name: "refresh_token", Value: "encodedRefreshToken"})
				mockTokenHandler.EXPECT().ValidateRefreshToken(gomock.Any()).Return(nil)
				mockTokenService.EXPECT().Find(gomock.Any()).Return(userId, nil)
				mockTokenHandler.EXPECT().CreateAccessToken(gomock.Any(), gomock.Any()).Return("", errors.New("internal token handler error"))
			},
			expectedCode:   500,
			wantedTokens:   []string{},
			wantedResponse: "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			isError:        true,
		},
		{
			title: "wrong refresh_token cookie and 403 response",
			args:  args{acessToken: "", userId: 0},
			mock: func(recorder *httptest.ResponseRecorder, userId int, accessToken string) {
				http.SetCookie(recorder, &http.Cookie{Name: "wrong_cookie", Value: "encodedRefreshToken"})
			},
			expectedCode:   403,
			wantedTokens:   []string{},
			wantedResponse: "\"http: named cookie not present\"",
			isError:        true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			router := gin.Default()
			router.GET("/refresh", authHandler.RefreshAccessToken)
			req, _ := http.NewRequest("GET", "/refresh", nil)
			recorder := httptest.NewRecorder()
			test.mock(recorder, test.args.userId, test.args.acessToken)
			req.Header = http.Header{"Cookie": recorder.Result().Header["Set-Cookie"]}
			router.ServeHTTP(recorder, req)
			if !test.isError {
				coockies := recorder.Result().Cookies()
				for _, c := range coockies {
					if c.Name == "access_token" {
						assert.Equal(t, test.wantedTokens[0], c.Value)
					} else if c.Name == "refresh_token" {
						assert.Equal(t, test.wantedTokens[1], c.Value)
					}
				}

			}
			assert.Equal(t, test.wantedResponse, recorder.Body.String())
			assert.Equal(t, test.expectedCode, recorder.Code)
		})
	}
}

func TestLogout(t *testing.T) {
	cntr := gomock.NewController(t)
	mockUserService := mocks.NewMockUserService(cntr)
	mockTokenHandler := mocks.NewMockTokenHandler(cntr)
	mockTokenService := mocks.NewMockTokenService(cntr)
	authHandler := NewAuthHandler(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, "localhost", 15, 15)
	type mockCall func() *http.Request
	testCases := []struct {
		title          string
		mock           mockCall
		wantedResponse string
		expectedCode   int
		wantedTokens   []string
		isError        bool
	}{
		{
			title: "success logout and 200 response",
			mock: func() *http.Request {
				mockTokenHandler.EXPECT().ValidateRefreshToken(gomock.Any()).Return(nil)
				mockTokenService.EXPECT().Remove(gomock.Any()).Return(nil)
				req, err := http.NewRequest("GET", "/logout", nil)
				assert.Nil(t, err)
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "encodedRefreshToken", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				req.AddCookie(&http.Cookie{Name: "access_token", Value: "encodedAccessToken", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				req.AddCookie(&http.Cookie{Name: "loggin", Value: "true", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				return req

			},
			expectedCode:   200,
			wantedTokens:   []string{"", ""},
			wantedResponse: "{\"status\":\"success\"}",
			isError:        false,
		},
		{
			title: "refresh token is missing and 403 response",
			mock: func() *http.Request {
				req, err := http.NewRequest("GET", "/logout", nil)
				assert.Nil(t, err)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: "encodedAccessToken", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				req.AddCookie(&http.Cookie{Name: "loggin", Value: "true", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				return req

			},
			expectedCode:   403,
			wantedResponse: "\"http: named cookie not present\"",
			isError:        true,
		},
		{
			title: "cannot delete refresh token  and 500 response",
			mock: func() *http.Request {
				req, err := http.NewRequest("GET", "/logout", nil)
				assert.Nil(t, err)
				mockTokenHandler.EXPECT().ValidateRefreshToken(gomock.Any()).Return(nil)
				mockTokenService.EXPECT().Remove(gomock.Any()).Return(errs.New(errs.Database))
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "NotValidRefreshToken", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				return req

			},
			expectedCode:   500,
			wantedResponse: "\"{\\\"error\\\":{\\\"kind\\\":\\\"internal_error\\\",\\\"message\\\":\\\"internal server error - please contact support\\\"}}\"",
			isError:        true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/logout", authHandler.Logout)
			req := test.mock()
			router.ServeHTTP(recorder, req)
			if !test.isError {
				coockies := recorder.Result().Cookies()
				for _, c := range coockies {
					if c.Name == "access_token" {
						assert.Equal(t, test.wantedTokens[0], c.Value)
					} else if c.Name == "refresh_token" {
						assert.Equal(t, test.wantedTokens[1], c.Value)
					}
				}

			}
			assert.Equal(t, test.wantedResponse, recorder.Body.String())
			assert.Equal(t, test.expectedCode, recorder.Code)

		})
	}
}
