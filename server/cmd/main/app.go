package main

import (
	"fmt"

	"github.com/VrMolodyakov/stock-market/internal"
	"github.com/VrMolodyakov/stock-market/internal/config"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("init app")

	cfg := config.GetConfig()
	logger := logging.GetLogger("debug")
	app := internal.NewApp(logger, cfg, gin.Default())
	app.Run()
}
