package route

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	SignUpUser(ctx *gin.Context)
	SignInUser(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authRouter struct {
	authHandler    AuthHandler
	authMiddleware AuthMiddleware
}

func NewAuthRouter(authHandler AuthHandler, authMiddleware AuthMiddleware) *authRouter {
	return &authRouter{authHandler: authHandler, authMiddleware: authMiddleware}
}

func (a *authRouter) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/register", a.authHandler.SignUpUser)
	router.POST("/login", a.authHandler.SignInUser)
	router.GET("/refresh", a.authHandler.RefreshAccessToken)
	router.GET("/logout", a.authMiddleware.Auth(), a.authHandler.Logout)
}
