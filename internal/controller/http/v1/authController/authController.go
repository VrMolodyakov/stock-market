package v1

import (
	"net/http"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/hashing"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/gin-gonic/gin"
)

type authController struct {
	logger       *logging.Logger
	userService  UserService
	tokenHandler TokenHandler
	tokenService TokenService
	accessTtl    int
	refreshTtl   int
}

func NewAuthController(
	userService UserService,
	logger *logging.Logger,
	tokenHandler TokenHandler,
	tokenService TokenService,
	accessTtl int,
	refreshTtl int) *authController {
	return &authController{
		userService:  userService,
		logger:       logger,
		tokenHandler: tokenHandler,
		tokenService: tokenService,
		accessTtl:    accessTtl,
		refreshTtl:   refreshTtl}
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

	ctx.SetCookie("access_token", accessToken, a.accessTtl*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, a.refreshTtl*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", a.accessTtl*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (a *authController) RefreshAccessToken(ctx *gin.Context) {
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
	ctx.SetCookie("access_token", accessToken, a.accessTtl*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", a.accessTtl*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (a *authController) Logout(ctx *gin.Context) {
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
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}
