package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/controller/http/v1/authController/mocks"
	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
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
	//testTime, _ := now.MarshalText()
	inputTime := now.Format(time.RFC3339) //string(testTime)
	authController := NewAuthController(mockUserService, logging.GetLogger("debug"), mockTokenHandler, mockTokenService, 15, 15)
	type mockCall func()
	testCases := []struct {
		title        string
		mock         mockCall
		inputRequest string
		want         string
	}{
		{
			title: "sign up and 201 response",
			mock: func() {
				user := entity.User{Username: "username", Password: "hashedpassword", CreateAt: now, Id: 1}
				mockUserService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
			},
			inputRequest: `{"username":"my_username","password":"my_password"}`,
			want:         fmt.Sprintf("{\"data\":{\"user\":{\"username\":\"username\",\"password\":\"hashedpassword\",\"create_at\":\"%v\"}},\"status\":\"success\"}", inputTime),
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
			assert.Equal(t, test.want, clearResponse(recorder.Body.String()))
			assert.Equal(t, 201, recorder.Code)
		})
	}
}

func clearResponse(s string) string {
	temp := strings.ReplaceAll(s, " ", "")
	return strings.ReplaceAll(temp, "\n", "")
}
