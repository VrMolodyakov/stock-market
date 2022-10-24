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
		title        string
		inputRequest string
		mockCall     mockCall
	}{
		{
			title:        "sign up and 201 response",
			inputRequest: `{}`,
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
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(response, nil)

			},
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			router := gin.Default()
			router.GET("/api/stock/symbols/TEST", stockHandler.GetStockInfo)
			req, _ := http.NewRequest("GET", "/api/stock/symbols/TEST", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
		})
	}
}
