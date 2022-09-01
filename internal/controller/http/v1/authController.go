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
	Create(ctx context.Context, username string, password string) (int, error)
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
	var user UserDto

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		errs.HTTPErrorResponse(ctx, a.logger, err)
		return
	}
	hashedPassword, err := hashing.HashPassword(user.Password)

}
