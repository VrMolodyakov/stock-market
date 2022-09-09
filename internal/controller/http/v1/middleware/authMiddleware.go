package middleware

import (
	"context"
	"strings"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetById(ctx context.Context, id int) (entity.User, error)
}

type TokenHandler interface {
	ValidateAccessToken(token string) (int, error)
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
		a.logger.Info("---------------INSIDE MIDDLEWARE---------------")
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
		userId, err := a.tokenHandler.ValidateAccessToken(accessToken)
		if err != nil {
			errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
			return
		}
		user, err := a.userService.GetById(ctx, userId)
		if err != nil {
			errs.HTTPErrorResponse(ctx, a.logger, err)
			return
		}
		a.logger.Info("SET CURRENT USER ")
		ctx.Set("user", user)
		ctx.Next()
	}

}
