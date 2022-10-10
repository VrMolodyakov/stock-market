package service

import (
	"time"

	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
)

type StockStorage interface {
	Set(symbol string, stockInfo string, expireAt time.Duration) error
	Get(symbol string) (string, error)
}

type cacheService struct {
	logger  *logging.Logger
	storage StockStorage
}

func NewCacheService(logger *logging.Logger, storage StockStorage) *cacheService {
	return &cacheService{logger: logger, storage: storage}
}

func (cs *cacheService) Save(symbol string, stockInfo string, duration time.Duration) error {
	if len(symbol) == 0 {
		return errs.New(errs.Validation, errs.Code("symbol is empty"), errs.Parameter("symbol"))
	}
	if len(stockInfo) == 0 {
		return errs.New(errs.Validation, errs.Code("stock info is empty"), errs.Parameter("stockInfo"))
	}
	return cs.storage.Set(symbol, stockInfo, duration)
}

func (cs *cacheService) Get(symbol string) (string, error) {
	cs.logger.Infof("try to get from cache symbol = %v", symbol)
	if len(symbol) == 0 {
		return "", errs.New(errs.Validation, errs.Code("symbol is empty"), errs.Parameter("symbol"))
	}
	return cs.storage.Get(symbol)
}
