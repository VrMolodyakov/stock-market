package middleware

import (
	"context"
	"strings"

	v1 "github.com/VrMolodyakov/stock-market/internal/controller/http/v1/auth"
	"github.com/VrMolodyakov/stock-market/internal/domain/entity"
	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetById(ctx context.Context, id int) (entity.User, error)
}

type TokenHandler interface {
	ValidateAccessToken(token string) (interface{}, error)
}

type TokenService interface {
	Find(refreshToken string) (int, error)
}

type authMiddleware struct {
	userService  UserService
	logger       *logging.Logger
	tokenHandler TokenHandler
	tokenService TokenService
}

func NewAuthMiddleware(userService UserService, tokenService TokenService, tokenHandler TokenHandler, logger *logging.Logger) *authMiddleware {
	return &authMiddleware{userService: userService, tokenService: tokenService, tokenHandler: tokenHandler, logger: logger}
}

func (a *authMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		coockie, err := ctx.Cookie("access_token")
		authHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authHeader)
		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = coockie
		}
		if accessToken == "" {
			errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
			return
		}
		sub, err := a.tokenHandler.ValidateAccessToken(accessToken)
		if err != nil {
			errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
			return
		}
		userId := sub.(float64)
		user, err := a.userService.GetById(ctx, int(userId))
		if err != nil {
			errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Validation, errs.Parameter("user id not found")))
			return
		}
		a.logger.Infof("SET CURRENT USER %v", user)
		ctx.Set("user", v1.User{Username: user.Username, CreateAt: user.CreateAt})
		ctx.Next()
	}

}
