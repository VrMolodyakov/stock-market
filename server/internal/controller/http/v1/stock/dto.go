package stock

type ChartResponse struct {
	Chart Chart `json:"chart"`
}

type Chart struct {
	Result []Result    `json:"result"`
	Error  interface{} `json:"error"`
}

type Result struct {
	Meta       Meta  `json:"meta"`
	Timestamp  []int `json:"timestamp"`
	Indicators struct {
		Quote []Quote `json:"quote"`
	} `json:"indicators"`
}

type Quote struct {
	Low    []float64 `json:"low"`
	Close  []float64 `json:"close"`
	Volume []float64 `json:"volume"`
	Open   []float64 `json:"open"`
	High   []float64 `json:"high"`
}

type Meta struct {
	Symbol             string  `json:"symbol"`
	RegularMarketTime  int     `json:"regularMarketTime"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}
