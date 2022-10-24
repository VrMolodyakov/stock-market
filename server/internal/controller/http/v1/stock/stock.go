package stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/VrMolodyakov/stock-market/pkg/metric"
	"github.com/gin-gonic/gin"
)

// gmt time

const expireAt int = 60
const closedTradeExpire int = 600

type CacheService interface {
	Save(symbol string, stockInfo string, duration time.Duration) error
	Get(symbol string) (string, error)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type stockService struct {
	metric       metric.Metric
	logger       *logging.Logger
	cacheService CacheService
	http         HttpClient
}

func NewStockHandler(metric metric.Metric, logger *logging.Logger, cache CacheService, http HttpClient) *stockService {
	return &stockService{metric: metric, logger: logger, cacheService: cache, http: http}
}

func (ss *stockService) GetStockInfo(ctx *gin.Context) {
	start := time.Now()
	code := ctx.Param("symbol")
	stockInfo, err := ss.cacheService.Get(code)
	if err == nil {
		ss.logger.Info("get symbol = %v from cache", code)
		b := []byte(stockInfo)
		var chart ChartResponse
		err = json.Unmarshal(b, &chart)
		if err != nil {
			ss.metric.HTTPResponseCounter.WithLabelValues(code, "500").Inc()
			errs.HTTPErrorResponse(ctx, ss.logger, errs.New(errs.Internal, err))
			return
		}
		ctx.JSON(http.StatusOK, chart)

	} else {
		url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%v", code)
		ss.logger.Info(url)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			ss.metric.HTTPResponseCounter.WithLabelValues(code, "500").Inc()
			errs.HTTPErrorResponse(ctx, ss.logger, errs.New(errs.Internal, err))
			return
		}
		resp, err := ss.http.Do(req)
		if err != nil {
			ss.metric.HTTPResponseCounter.WithLabelValues(code, "500").Inc()
			errs.HTTPErrorResponse(ctx, ss.logger, errs.New(errs.Internal, err))
			return
		}
		defer resp.Body.Close()
		var chart ChartResponse
		err = json.NewDecoder(resp.Body).Decode(&chart)

		if err != nil {
			ss.metric.HTTPResponseCounter.WithLabelValues(code, "500").Inc()
			errs.HTTPErrorResponse(ctx, ss.logger, errs.New(errs.Internal, err))
			return
		}
		ss.cache(chart, code)
		dur := float64(time.Since(start).Milliseconds())
		ss.metric.ResponseDurationHistogram.WithLabelValues(code).Observe(dur)
		ss.metric.HTTPResponseCounter.WithLabelValues(code, "200").Inc()
		ctx.JSON(http.StatusOK, chart)
	}

}

func (ss *stockService) cache(chart ChartResponse, symbol string) {
	ss.logger.Infof("try to save in cache symbol = %v", symbol)
	stringPayload, _ := json.Marshal(chart)
	h := time.Now().UTC().Hour()
	var err error
	if h >= 9 || h <= 16 {
		err = ss.cacheService.Save(symbol, string(stringPayload), time.Duration(expireAt)*time.Second)
	} else {
		err = ss.cacheService.Save(symbol, string(stringPayload), time.Duration(closedTradeExpire)*time.Second)
	}

	if err != nil {
		ss.logger.Errorf("cannot save to cache due to : %v", err)
	}
}
