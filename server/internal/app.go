package internal

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	stockstorage "github.com/VrMolodyakov/stock-market/internal/adapter/stockStorage"
	"github.com/VrMolodyakov/stock-market/internal/adapter/tokenStorage"
	userstorage "github.com/VrMolodyakov/stock-market/internal/adapter/userStorage"
	"github.com/VrMolodyakov/stock-market/internal/config"
	v1 "github.com/VrMolodyakov/stock-market/internal/controller/http/v1/auth"
	"github.com/VrMolodyakov/stock-market/internal/controller/http/v1/middleware"
	"github.com/VrMolodyakov/stock-market/internal/controller/http/v1/route"
	"github.com/VrMolodyakov/stock-market/internal/controller/http/v1/stock"
	userController "github.com/VrMolodyakov/stock-market/internal/controller/http/v1/user"
	"github.com/VrMolodyakov/stock-market/internal/domain/service"
	"github.com/VrMolodyakov/stock-market/pkg/client/postgresql"
	"github.com/VrMolodyakov/stock-market/pkg/client/redis"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/VrMolodyakov/stock-market/pkg/shutdown"
	"github.com/VrMolodyakov/stock-market/pkg/token"
	"github.com/gin-gonic/gin"
)

const (
	writeTimeout = 15 * time.Second
	readTimeout  = 15 * time.Second
)

type app struct {
	logger *logging.Logger
	cfg    *config.Config
	server *gin.Engine
}

func NewApp(logger *logging.Logger, cfg *config.Config, server *gin.Engine) *app {
	return &app{cfg: cfg, logger: logger, server: server}
}

func (a *app) Run() {
	a.startHttp()
}

func (a *app) startHttp() {
	a.logger.Info("start http server")
	pgCfg := postgresql.NewPgConfig(
		a.cfg.PostgreSql.Username,
		a.cfg.PostgreSql.Password,
		a.cfg.PostgreSql.Host,
		a.cfg.PostgreSql.Port,
		a.cfg.PostgreSql.Dbname,
		a.cfg.PostgreSql.PoolSize)
	psqlClient, _ := postgresql.NewClient(context.Background(), 5, 5*time.Second, pgCfg)
	storage := userstorage.New(a.logger, psqlClient)
	rdCfg := redis.NewRdConfig(a.cfg.Redis.Password, a.cfg.Redis.Host, a.cfg.Redis.Port, a.cfg.Redis.DbNumber)
	rdClient, err := redis.NewClient(context.Background(), &rdCfg)

	if err != nil {
		a.logger.Fatal(err)
	}
	stockStorage := stockstorage.NewStockStorage(a.logger, rdClient)
	tokenStorage := tokenStorage.NewChoiceCache(rdClient, a.logger)
	accessPair, refreshPair := a.initTokens()
	tokenHandler := token.NewTokenHandler(a.logger, accessPair, refreshPair)
	tokenService := service.NewTokenService(tokenStorage, a.logger)
	userService := service.NewUserService(a.logger, storage)
	cacheService := service.NewCacheService(a.logger, stockStorage)
	authController := v1.NewAuthController(userService, a.logger, tokenHandler, tokenService, a.cfg.Token.AccessTtl, a.cfg.Token.RefreshTtl)
	authMiddleware := middleware.NewAuthMiddleware(userService, tokenService, tokenHandler, a.logger)
	userController := userController.NewUserController(userService, a.logger)
	stockHandler := stock.NewStockHandler(*a.logger, cacheService)
	a.server.Use(middleware.CORSMiddleware())
	router := a.server.Group("/api")
	authRouter := route.NewAuthRouter(authController, authMiddleware)
	userRouter := route.NewUserRouter(userController, authMiddleware)
	stockRouter := route.NewStockRouter(stockHandler, authMiddleware)
	authRouter.AuthRoute(router)
	userRouter.UserRoute(router)
	stockRouter.StockRoute(router)

	a.server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": fmt.Sprintf("Route %s not found", ctx.Request.URL)})
	})

	port := fmt.Sprintf(":%s", a.cfg.Port)
	server := &http.Server{
		Addr:         port,
		Handler:      a.server,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, rdClient, server)
	defer psqlClient.Close()
	if err := server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	a.logger.Info("app shutdown")

}

func (a *app) initTokens() (token.TokenPair, token.TokenPair) {
	aprk, err := base64.StdEncoding.DecodeString(a.cfg.Token.AccessPrivate)
	a.checkErr(err)
	apbk, err := base64.StdEncoding.DecodeString(a.cfg.Token.AccessPublic)
	a.checkErr(err)
	rprk, err := base64.StdEncoding.DecodeString(a.cfg.Token.RefreshPrivate)
	a.checkErr(err)
	rpbk, err := base64.StdEncoding.DecodeString(a.cfg.Token.RefreshPublic)
	a.checkErr(err)
	apair := token.TokenPair{PrivateKey: aprk, PublicKey: apbk}
	rpair := token.TokenPair{PrivateKey: rprk, PublicKey: rpbk}
	return apair, rpair
}

func (a *app) checkErr(err error) {
	if err != nil {
		a.logger.Fatal(err)
	}
}
