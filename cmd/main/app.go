package main

import (
	"fmt"

	"github.com/VrMolodyakov/jwt-auth/internal"
	"github.com/VrMolodyakov/jwt-auth/internal/config"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("init app")

	cfg := config.GetConfig()
	logger := logging.GetLogger("debug")
	app := internal.NewApp(logger, cfg, gin.Default())
	app.Run()
}
