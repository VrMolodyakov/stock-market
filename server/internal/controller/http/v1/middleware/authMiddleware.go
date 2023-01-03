package middleware

import (
	"context"
	"strings"

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
			ctx.Abort()
			errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
			return
		}
		a.logger.Debug("access token:", accessToken)
		sub, err := a.tokenHandler.ValidateAccessToken(accessToken)
		if err != nil {
			ctx.Abort()
			errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
			return
		}
		userId := sub.(float64)
		user, err := a.userService.GetById(ctx, int(userId))
		if err != nil {
			ctx.Abort()
			errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Validation, errs.Parameter("user id not found")))
			return
		}
		a.logger.Debugf("set current context user %v = ", user)
		ctx.Set("user", user)
		ctx.Next()
	}

}
