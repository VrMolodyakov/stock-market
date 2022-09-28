package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/VrMolodyakov/stock-market/pkg/logging"
)

const requestTime int = 3

type stockService struct {
	logger logging.Logger
}

func NewStockService(logger logging.Logger) *stockService {
	return &stockService{logger: logger}
}

func (ss *stockService) GetStockInfo(code string) {
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*time.Duration(requestTime))
	defer cncl()
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%v", code)
	ss.logger.Info(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {

	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}
