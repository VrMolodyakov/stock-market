package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/adapter/tokenStorage"
	userstorage "github.com/VrMolodyakov/jwt-auth/internal/adapter/userStorage"
	"github.com/VrMolodyakov/jwt-auth/internal/config"
	v1 "github.com/VrMolodyakov/jwt-auth/internal/controller/http/v1"
	"github.com/VrMolodyakov/jwt-auth/internal/domain/service"
	"github.com/VrMolodyakov/jwt-auth/pkg/client/postgresql"
	"github.com/VrMolodyakov/jwt-auth/pkg/client/redis"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/VrMolodyakov/jwt-auth/pkg/token"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("init app")

	cfg := config.GetConfig()
	logger := logging.GetLogger("debug")
	pgCfg := postgresql.NewPgConfig(
		cfg.PostgreSql.Username,
		cfg.PostgreSql.Password,
		cfg.PostgreSql.Host,
		cfg.PostgreSql.Port,
		cfg.PostgreSql.Dbname,
		cfg.PostgreSql.PoolSize)
	client, _ := postgresql.NewClient(context.Background(), 5, 5*time.Second, pgCfg)
	storage := userstorage.New(logger, client)
	rdCfg := redis.NewRdConfig(cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DbNumber)
	rdClient, err := redis.NewClient(context.Background(), &rdCfg)
	if err != nil {
		logger.Fatal(err)
	}
	tokenStorage := tokenStorage.NewChoiceCache(rdClient, logger)
	aprk, _ := base64.StdEncoding.DecodeString(cfg.Token.AccessPrivate)
	apbk, _ := base64.StdEncoding.DecodeString(cfg.Token.AccessPublic)
	rprk, _ := base64.StdEncoding.DecodeString(cfg.Token.RefreshPrivate)
	rpbk, _ := base64.StdEncoding.DecodeString(cfg.Token.RefreshPublic)
	apair := token.TokenPair{PrivateKey: aprk, PublicKey: apbk}
	rpair := token.TokenPair{PrivateKey: rprk, PublicKey: rpbk}
	tokenHandler := token.NewTokenHandler(logger, apair, rpair)
	tokenService := service.NewTokenService(tokenStorage, logger)
	serv := service.NewUserStorage(logger, storage)
	cntrl := v1.NewAuthController(serv, logger, tokenHandler, tokenService, 15, 15)
	r := gin.Default()
	r.POST("/user", cntrl.SignUpUser)
	r.POST("/create", cntrl.SignInUser)
	r.POST("/refresh", cntrl.RefreshAccessToken)
	r.Run()
}
