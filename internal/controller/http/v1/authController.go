package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/hashing"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Create(ctx context.Context, username string, password string) (entity.User, error)
	Get(ctx context.Context, username string) (entity.User, error)
}

type TokenHandler interface {
	CreateAccessToken(ttl time.Duration, payload interface{}) (string, error)
	CreateRefreshToken(ttl time.Duration, payload interface{}) (string, error)
	ValidateAccessToken(token string) (interface{}, error)
	ValidateRefreshToken(token string) (interface{}, error)
}

type authController struct {
	logger       *logging.Logger
	userService  UserService
	tokenHandler TokenHandler
	accessTtl    int
	refreshTtl   int
}

func NewAuthController(userService UserService, logger *logging.Logger, tokenHandler TokenHandler, accessTtl int, refreshTtl int) *authController {
	return &authController{userService: userService, logger: logger, tokenHandler: tokenHandler, accessTtl: accessTtl, refreshTtl: refreshTtl}
}

func (a *authController) SignUpUser(ctx *gin.Context) {
	var request UserRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	hashedPassword, err := hashing.HashPassword(request.Password)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	user, err := a.userService.Create(ctx, request.Username, hashedPassword)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	response := ResponseFromEntity(user)
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": response}})

}

func (a *authController) SignInUser(ctx *gin.Context) {
	var request UserRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}

	user, err := a.userService.Get(ctx, request.Username)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Validation, errs.Code("wrong username"), errs.Parameter("username"), err))
		return
	}
	a.logger.Debugf("FIND USER = %v", user)
	err = hashing.ComparePassword(user.Password, request.Password)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Validation, errs.Code("wrong password"), errs.Parameter("password"), err))
		return
	}

	accessToken, err := a.tokenHandler.CreateAccessToken(time.Duration(a.accessTtl)*time.Second, user.Id)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Internal, err))
		return
	}
	refreshToken, err := a.tokenHandler.CreateRefreshToken(time.Duration(a.refreshTtl)*time.Second, user.Id)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Validation, errs.Code("wrong password"), errs.Parameter("password"), err))
		return
	}
	ctx.SetCookie("access_token", accessToken, a.accessTtl*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, a.refreshTtl*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", a.accessTtl*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}
