package v1

import (
	"context"

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

type authController struct {
	logger      *logging.Logger
	userService UserService
}

func NewAuthController(userService UserService, logger *logging.Logger) *authController {
	return &authController{userService: userService, logger: logger}
}

func (a *authController) SignUpUser(ctx *gin.Context) {
	var dto UserDto

	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	hashedPassword, err := hashing.HashPassword(dto.Password)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	user, err := a.userService.Create(ctx, dto.Username, hashedPassword)

}
