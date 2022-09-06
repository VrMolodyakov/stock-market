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

func (t *tokenStorage) Set(refreshToken string, userId int, expireAt time.Duration) error {
	t.logger.Debugf("try to save token = %v for user with id = %v", refreshToken, userId)
	err := t.client.Set(refreshToken, strconv.Itoa(userId), expireAt).Err()
	if err != nil {
		return errs.New(errs.Database, err)
	}
	return nil
}

func (t *tokenStorage) Get(refreshToken string) (int, error) {
	value, err := t.client.Get(refreshToken).Result()
	if err != nil {
		return -1, errs.New(errs.Database, err)
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		return -1, errs.New(errs.Validation, errs.Code("couldn't parse user id"), errs.Parameter("user_id"), err)
	}
	return count, nil
}

func (t *tokenStorage) Delete(refreshToken string) error {
	err := t.client.Del(refreshToken).Err()
	if err != nil {
		return errs.New(errs.Database, err)
	}
	return nil
}
