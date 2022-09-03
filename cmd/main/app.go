package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/config"
	"github.com/VrMolodyakov/jwt-auth/internal/domain/entity"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/VrMolodyakov/jwt-auth/pkg/token"
)

func main() {
	fmt.Println("init app")

	cfg := config.GetConfig()
	logger := logging.GetLogger("debug")
	// pgCfg := postgresql.NewPgConfig(
	// 	cfg.PostgreSql.Username,
	// 	cfg.PostgreSql.Password,
	// 	cfg.PostgreSql.Host,
	// 	cfg.PostgreSql.Port,
	// 	cfg.PostgreSql.Dbname,
	// 	cfg.PostgreSql.PoolSize)
	// client, _ := postgresql.NewClient(context.Background(), 5, 5*time.Second, pgCfg)
	// storage := userstorage.New(logger, client)
	// aprk, _ := base64.StdEncoding.DecodeString(cfg.Token.AccessPrivate)
	// apbk, _ := base64.StdEncoding.DecodeString(cfg.Token.AccessPublic)
	// rprk, _ := base64.StdEncoding.DecodeString(cfg.Token.RefreshPrivate)
	// rpbk, _ := base64.StdEncoding.DecodeString(cfg.Token.RefreshPublic)
	// apair := token.TokenPair{PrivateKey: aprk, PublicKey: apbk}
	// rpair := token.TokenPair{PrivateKey: rprk, PublicKey: rpbk}
	// tokenHandler := token.NewTokenHandler(logger, apair, rpair)
	// serv := service.NewUserStorage(logger, storage)
	// cntrl := v1.NewAuthController(serv, logger, tokenHandler, 15, 15)
	// r := gin.Default()
	// r.POST("/user", cntrl.SignUpUser)
	// r.POST("/create", cntrl.SignInUser)
	// r.Run()

	aprk, _ := base64.StdEncoding.DecodeString(cfg.Token.AccessPrivate)
	apbk, _ := base64.StdEncoding.DecodeString(cfg.Token.AccessPublic)
	rprk, _ := base64.StdEncoding.DecodeString(cfg.Token.RefreshPrivate)
	rpbk, _ := base64.StdEncoding.DecodeString(cfg.Token.RefreshPublic)

	apair := token.TokenPair{PrivateKey: aprk, PublicKey: apbk}
	rpair := token.TokenPair{PrivateKey: rprk, PublicKey: rpbk}
	tokenHandler := token.NewTokenHandler(logger, apair, rpair)
	tokenHandler.CreateAccessToken(5*time.Second, entity.User{Id: 1})
}
