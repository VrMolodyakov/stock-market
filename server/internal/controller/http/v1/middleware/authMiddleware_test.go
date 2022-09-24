package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/VrMolodyakov/stock-market/internal/controller/http/v1/auth"
	"github.com/VrMolodyakov/stock-market/internal/controller/http/v1/middleware/mocks"
	"github.com/VrMolodyakov/stock-market/internal/domain/entity"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	cntr := gomock.NewController(t)
	tokenHandler := mocks.NewMockTokenHandler(cntr)
	tokenService := mocks.NewMockTokenService(cntr)
	userService := mocks.NewMockUserService(cntr)
	logger := logging.GetLogger("debug")
	authMiddleware := NewAuthMiddleware(userService, tokenService, tokenHandler, logger)
	type mockCall func(req *http.Request)
	testCases := []struct {
		title      string
		mockCall   mockCall
		handler    func(ctx *gin.Context)
		wantedCode int
		wantedBody string
	}{
		{
			title: "find access token and succes response",
			mockCall: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "access_token", Value: "encodedAccessToken", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				var i interface{}
				i = 1.0
				user := entity.User{Id: 1, Username: "some-username"}
				tokenHandler.EXPECT().ValidateAccessToken(gomock.Any()).Return(i, nil)
				userService.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			handler: func(ctx *gin.Context) {
				user := ctx.MustGet("user").(v1.User)
				assert.Equal(t, "some-username", user.Username)
				ctx.JSON(http.StatusOK, "success")
			},
			wantedCode: 200,
			wantedBody: "\"success\"",
		},
		{
			title: "cannot find access token and 403 response",
			mockCall: func(req *http.Request) {
			},
			handler: func(ctx *gin.Context) {
			},
			wantedCode: 403,
			wantedBody: "\"http: named cookie not present\"",
		},
		{
			title: "cannot validate access token and 403 response",
			mockCall: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "access_token", Value: "encodedAccessToken", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				tokenHandler.EXPECT().ValidateAccessToken(gomock.Any()).Return(nil, errors.New("token handler internal error"))

			},
			handler: func(ctx *gin.Context) {
			},
			wantedCode: 403,
			wantedBody: "\"token handler internal error\"",
		},
		{
			title: "cannot find user id and 403 response",
			mockCall: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "access_token", Value: "encodedAccessToken", MaxAge: 60 * 60, Path: "/", Domain: "localhost", Secure: false, HttpOnly: true})
				var i interface{}
				i = 1.0
				tokenHandler.EXPECT().ValidateAccessToken(gomock.Any()).Return(i, nil)
				userService.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(entity.User{}, errors.New("cannnot find user with given id"))
			},
			handler: func(ctx *gin.Context) {
			},
			wantedCode: 400,
			wantedBody: "\"user id not found\"",
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			test.mockCall(req)
			assert.NoError(t, err)
			router.Use(authMiddleware.Auth())
			router.GET("/", test.handler)
			router.ServeHTTP(w, req)
			assert.Equal(t, test.wantedCode, w.Code)
			assert.Equal(t, test.wantedBody, w.Body.String())
		})
	}
}
