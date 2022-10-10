package stockstorage

import (
	"time"

	"github.com/VrMolodyakov/stock-market/internal/errs"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/go-redis/redis"
)

type stockCache struct {
	logger *logging.Logger
	client *redis.Client
}

func NewStockStorage(logger *logging.Logger, client *redis.Client) *stockCache {
	return &stockCache{logger: logger, client: client}
}

func (sc *stockCache) Set(symbol string, stockInfo string, expireAt time.Duration) error {
	sc.logger.Debugf("try to save for symbol = %v", symbol)
	err := sc.client.Set(symbol, stockInfo, expireAt).Err()
	if err != nil {
		return errs.New(errs.Database, err)
	}
	return nil
}

func (sc *stockCache) Get(symbol string) (string, error) {
	sc.logger.Infof("try to get %v", symbol)
	url, err := sc.client.Get(symbol).Result()
	if err != nil {
		sc.logger.Errorf("cannot get full url for short url = %v due to ", err)
		return "", errs.New(errs.Database, err)
	}
	return url, err
}
