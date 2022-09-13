package v1

import (
	"bytes"
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
