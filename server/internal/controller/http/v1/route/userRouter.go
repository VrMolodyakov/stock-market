package route

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetCurrentUser(ctx *gin.Context)
}

type userRouter struct {
	userHandler    UserHandler
	authMiddleware AuthMiddleware
}

func NewUserRouter(userHandler UserHandler, authMiddleware AuthMiddleware) *userRouter {
	return &userRouter{userHandler: userHandler, authMiddleware: authMiddleware}
}

func (a *userRouter) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/users")
	router.GET("/me", a.authMiddleware.Auth(), a.userHandler.GetCurrentUser)
}
