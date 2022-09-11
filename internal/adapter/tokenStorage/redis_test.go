package tokenStorage

import (
	"errors"
	"testing"
	"time"

	"github.com/VrMolodyakov/jwt-auth/pkg/logging"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var (
	redisServer *miniredis.Miniredis
	redisClient *redis.Client
)

func TestSet(t *testing.T) {
	setUp()
	defer teardown()
	repo := NewChoiceCache(redisClient, logging.GetLogger("debug"))
	type args struct {
		refreshToken string
		userId       int
		expire       time.Duration
	}
	type mockCall func()
	testCases := []struct {
		title   string
		input   args
		isError bool
		mock    mockCall
	}{
		{
			title:   "Success Set()",
			input:   args{refreshToken: "vote", userId: 1, expire: 1 * time.Second},
			mock:    func() {},
			isError: false,
		},
		{
			title: "reddis internal error and Set() return error",
			input: args{refreshToken: "vote", userId: 1, expire: 1 * time.Second},
			mock: func() {
				redisServer.SetError("interanl redis error")
			},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			err := repo.Set(test.input.refreshToken, test.input.userId, test.input.expire)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	setUp()
	defer teardown()

	repo := NewChoiceCache(redisClient, logging.GetLogger("debug"))
	type mockCall func(refreshToken string, userId int, expire time.Duration) error
	type args struct {
		refreshToken string
		userId       int
		expire       time.Duration
	}
	testCases := []struct {
		title   string
		input   args
		isError bool
		mock    mockCall
		want    int
	}{
		{
			title:   "Get should find title and return count",
			input:   args{refreshToken: "refresh token", userId: 1, expire: 5 * time.Second},
			isError: false,
			mock: func(refreshToken string, userId int, expire time.Duration) error {
				return repo.Set(refreshToken, userId, expire)
			},
			want: 1,
		},
		{
			title:   "Get doens't find key and should return error",
			input:   args{refreshToken: "unsaved token", userId: 1, expire: 5 * time.Second},
			isError: true,
			mock: func(refreshToken string, userId int, expire time.Duration) error {
				return nil
			},
			want: -1,
		},
		{
			title:   "reddis internal error and Get return error ",
			input:   args{refreshToken: "refresh token", userId: 1, expire: 5 * time.Second},
			isError: true,
			mock: func(refreshToken string, userId int, expire time.Duration) error {
				redisServer.SetError("interanl redis error")
				return errors.New("internal error")
			},
			want: -1,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			_ = test.mock(test.input.refreshToken, test.input.userId, test.input.expire)
			got, err := repo.Get(test.input.refreshToken)
			if !test.isError {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			} else {
				assert.Error(t, err)
			}

		})
	}

}

func TestDelete(t *testing.T) {
	setUp()
	defer teardown()

	repo := NewChoiceCache(redisClient, logging.GetLogger("debug"))
	type mockCall func(refreshToken string, userId int, expire time.Duration) error
	type args struct {
		refreshToken string
		userId       int
		expire       time.Duration
	}
	testCases := []struct {
		title   string
		input   args
		isError bool
		mock    mockCall
	}{
		{
			title:   "Get should find title and return count",
			input:   args{refreshToken: "refresh token", userId: 1, expire: 5 * time.Second},
			isError: false,
			mock: func(refreshToken string, userId int, expire time.Duration) error {
				return repo.Set(refreshToken, userId, expire)
			},
		},
		{
			title:   "reddis internal error and Get return error ",
			input:   args{refreshToken: "refresh token", userId: 1, expire: 5 * time.Second},
			isError: true,
			mock: func(refreshToken string, userId int, expire time.Duration) error {
				redisServer.SetError("interanl redis error")
				return errors.New("internal error")
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			_ = test.mock(test.input.refreshToken, test.input.userId, test.input.expire)
			err := repo.Delete(test.input.refreshToken)
			if !test.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}

}

func setUp() {
	redisServer = mockRedis()
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func teardown() {
	redisServer.Close()
}
