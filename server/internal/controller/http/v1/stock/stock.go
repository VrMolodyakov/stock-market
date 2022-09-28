package stock

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/gin-gonic/gin"
)

type stockService struct {
	logger logging.Logger
}

func NewStockHandler(logger logging.Logger) *stockService {
	return &stockService{logger: logger}
}

func (ss *stockService) GetStockInfo(ctx *gin.Context) {
	code := ctx.Param("symbol")
	ss.logger.Info("code", code)
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%v", code)
	ss.logger.Info(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errs.HTTPErrorResponse(ctx, &ss.logger, errs.New(errs.Internal, err))
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errs.HTTPErrorResponse(ctx, &ss.logger, errs.New(errs.Internal, err))
		return
	}
	defer resp.Body.Close()
	var chart ChartResponse
	err = json.NewDecoder(resp.Body).Decode(&chart)
	if err != nil {
		errs.HTTPErrorResponse(ctx, &ss.logger, errs.New(errs.Internal, err))
		return
	}
	ctx.JSON(http.StatusOK, chart)
}
