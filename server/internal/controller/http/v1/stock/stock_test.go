package stock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VrMolodyakov/stock-market/internal/controller/http/v1/stock/mocks"
	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/VrMolodyakov/stock-market/pkg/metric"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignInUser(t *testing.T) {
	cntr := gomock.NewController(t)
	mockCacheService := mocks.NewMockCacheService(cntr)
	mockHttpClient := mocks.NewMockHttpClient(cntr)
	prometheusClient := metric.NewPrometheusClient(true)
	metric := metric.NewMetric(prometheusClient.Registry())
	stockHandler := NewStockHandler(metric, logging.GetLogger("debug"), mockCacheService, mockHttpClient)
	type mockCall func()
	testCases := []struct {
		title              string
		mockCall           mockCall
		expectedCode       int
		expectdSymbol      string
		expectedMarketTime int
		ExpectdMarketPrice float64
		isError            bool
	}{
		{
			title: "successful receipt of stock info and 200 response",
			mockCall: func() {
				chart := ChartResponse{
					Chart: Chart{
						Result: []Result{
							{
								Meta: Meta{
									Symbol:             "TEST",
									RegularMarketTime:  42,
									RegularMarketPrice: 42.0,
								},
							},
						},
						Error: nil,
					},
				}
				b, err := json.Marshal(chart)
				if err != nil {
					t.Fatal(err)
				}
				response := &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(string(b))),
				}
				mockCacheService.EXPECT().Get(gomock.Any()).Return("", errors.New("cache is empty"))
				mockCacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(response, nil)

			},
			expectedCode:       200,
			expectdSymbol:      "TEST",
			expectedMarketTime: 42,
			ExpectdMarketPrice: 42.0,
			isError:            false,
		},

		{
			title: "couldn't complete the request and 500 response",
			mockCall: func() {

				mockCacheService.EXPECT().Get(gomock.Any()).Return("", errors.New("cache is empty"))
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("http client error"))

			},
			expectedCode: 500,
			isError:      true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mockCall()
			router := gin.Default()
			router.GET("/api/stock/symbols/TEST", stockHandler.GetStockInfo)
			req, _ := http.NewRequest("GET", "/api/stock/symbols/TEST", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			if !test.isError {

				var chart ChartResponse
				err := json.NewDecoder(recorder.Body).Decode(&chart)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, test.expectdSymbol, chart.Chart.Result[0].Meta.Symbol)
				assert.Equal(t, test.expectedMarketTime, chart.Chart.Result[0].Meta.RegularMarketTime)
				assert.Equal(t, test.ExpectdMarketPrice, chart.Chart.Result[0].Meta.RegularMarketPrice)
			}
			assert.Equal(t, test.expectedCode, recorder.Code)
		})
	}
}
