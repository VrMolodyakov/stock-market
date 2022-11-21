package service

import (
	"testing"
	"time"

	"github.com/VrMolodyakov/stock-market/internal/domain/service/mocks"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCahceSave(t *testing.T) {
	cntr := gomock.NewController(t)
	stockStorage := mocks.NewMockStockStorage(cntr)
	cacheService := NewCacheService(logging.GetLogger("debug"), stockStorage)
	type mockCall func()
	type input struct {
		symbol    string
		stockInfo string
		duration  time.Duration
	}
	testCases := []struct {
		title    string
		mockCall mockCall
		input    input
		isError  bool
	}{
		{
			title: "success cache save",
			mockCall: func() {
				stockStorage.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			input:   input{symbol: "test", stockInfo: "testInfo", duration: 1 * time.Minute},
			isError: false,
		},
		{
			title: "empty symbol and return error",
			mockCall: func() {
			},
			input:   input{symbol: "", stockInfo: "testInfo", duration: 1 * time.Minute},
			isError: true,
		},
		{
			title: "empty stock info and return error",
			mockCall: func() {
			},
			input:   input{symbol: "test", stockInfo: "", duration: 1 * time.Minute},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mockCall()
			err := cacheService.Save(test.input.symbol, test.input.stockInfo, test.input.duration)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestCahceGet(t *testing.T) {
	cntr := gomock.NewController(t)
	stockStorage := mocks.NewMockStockStorage(cntr)
	cacheService := NewCacheService(logging.GetLogger("debug"), stockStorage)
	type mockCall func()

	testCases := []struct {
		title    string
		mockCall mockCall
		input    string
		isError  bool
		want     string
	}{
		{
			title: "success get from cache",
			mockCall: func() {
				stockStorage.EXPECT().Get(gomock.Any()).Return("test info", nil)
			},
			input:   "test symbol",
			want:    "test info",
			isError: false,
		},
		{
			title: "empty symbol and return error",
			mockCall: func() {
			},
			input:   "",
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mockCall()
			expected, err := cacheService.Get(test.input)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, expected, test.want)
				assert.NoError(t, err)
			}

		})
	}
}
