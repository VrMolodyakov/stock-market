package route

import "github.com/gin-gonic/gin"

type StockHandler interface {
	GetStockInfo(ctx *gin.Context)
}

type stockRouter struct {
	stockHandler   StockHandler
	authMiddleware AuthMiddleware
}

func NewStockRouter(stockHandler StockHandler, authMiddleware AuthMiddleware) *stockRouter {
	return &stockRouter{stockHandler: stockHandler, authMiddleware: authMiddleware}
}

func (s *stockRouter) StockRoute(rg *gin.RouterGroup) {
	router := rg.Group("/stock/symbols")
	router.GET("/:symbol", s.authMiddleware.Auth(), s.stockHandler.GetStockInfo)
}
