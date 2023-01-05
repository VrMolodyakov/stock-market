package v1

import (
	"net/http"
	"time"

	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/hashing"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	logger       *logging.Logger
	userService  UserService
	tokenHandler TokenHandler
	tokenService TokenService
	host         string
	accessTtl    int
	refreshTtl   int
}

func NewAuthHandler(
	userService UserService,
	logger *logging.Logger,
	tokenHandler TokenHandler,
	tokenService TokenService,
	host string,
	accessTtl int,
	refreshTtl int) *authHandler {
	return &authHandler{
		userService:  userService,
		logger:       logger,
		tokenHandler: tokenHandler,
		tokenService: tokenService,
		host:         host,
		accessTtl:    accessTtl,
		refreshTtl:   refreshTtl}
}

func (a *authHandler) SignUpUser(ctx *gin.Context) {
	var request UserRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Validation, errs.Code("incorrect data format")))
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

func (a *authHandler) SignInUser(ctx *gin.Context) {
	var request UserRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}

	user, err := a.userService.Get(ctx, request.Username)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	a.logger.Debugf("find user = %v", user)
	err = hashing.ComparePassword(user.Password, request.Password)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Validation, errs.Code("wrong password"), errs.Parameter("password"), err))
		return
	}

	accessToken, err := a.tokenHandler.CreateAccessToken(time.Duration(a.accessTtl)*time.Minute, user.Id)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Internal, err))
		return
	}
	refreshToken, err := a.tokenHandler.CreateRefreshToken(time.Duration(a.refreshTtl)*time.Minute, user.Id)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Internal, err))
		return
	}

	err = a.tokenService.Save(refreshToken, user.Id, time.Duration(a.refreshTtl)*time.Minute)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}

	ctx.SetCookie("access_token", accessToken, a.accessTtl*60, "/", a.host, false, true)
	ctx.SetCookie("refresh_token", refreshToken, a.refreshTtl*60, "/", a.host, false, true)
	ctx.SetCookie("logged_in", "true", a.accessTtl*60, "/", a.host, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (a *authHandler) RefreshAccessToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
		return
	}
	a.logger.Infof("get refreshToken from cookie = %v", refreshToken)
	err = a.tokenHandler.ValidateRefreshToken(refreshToken)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
		return
	}
	userId, err := a.tokenService.Find(refreshToken)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	accessToken, err := a.tokenHandler.CreateAccessToken(time.Duration(a.refreshTtl)*time.Minute, userId)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Internal, err))
		return
	}
	ctx.SetCookie("access_token", accessToken, a.accessTtl*60, "/", a.host, false, true)
	ctx.SetCookie("logged_in", "true", a.accessTtl*60, "/", a.host, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (a *authHandler) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
		return
	}
	err = a.tokenHandler.ValidateRefreshToken(refreshToken)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, errs.New(errs.Unauthorized, err))
		return
	}
	err = a.tokenService.Remove(refreshToken)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	ctx.SetCookie("access_token", "", -1, "/", a.host, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", a.host, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", a.host, false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}
