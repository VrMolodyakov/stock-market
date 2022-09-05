package tokenStorage

import (
	"strconv"
	"time"

	"github.com/VrMolodyakov/jwt-auth/internal/errs"
	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/go-redis/redis"
)

type tokenStorage struct {
	logger *logging.Logger
	client *redis.Client
}

func NewChoiceCache(client *redis.Client, logger *logging.Logger) *tokenStorage {
	return &tokenStorage{logger: logger, client: client}
}

func (c *tokenStorage) Set(refreshToken string, userId int, expireAt time.Duration) error {
	c.logger.Debugf("try to save token = %v for user with id = %v", refreshToken, userId)
	err := c.client.Set(refreshToken, strconv.Itoa(userId), expireAt).Err()
	if err != nil {
		return errs.New(errs.Database, err)
	}
	return nil
}

func (c *tokenStorage) Get(refreshToken string) (int, error) {
	value, err := c.client.Get(refreshToken).Result()
	if err != nil {
		return -1, errs.New(errs.Database, err)
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		return -1, errs.New(errs.Validation, errs.Code("couldn't parse user id"), errs.Parameter("user_id"), err)
	}
	return count, nil
}
