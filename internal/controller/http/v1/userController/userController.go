package userController

import (
	"context"
	"net/http"

	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Get(ctx context.Context, username string) (entity.User, error)
}

type userController struct {
	userService UserService
	logger      *logging.Logger
}

func NewUserController(userService UserService, logger *logging.Logger) *userController {
	return &userController{logger: logger, userService: userService}
}

func (u *userController) GetCurrentUser(ctx *gin.Context) {
	userId := ctx.MustGet("user").(int)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userId}})
}
