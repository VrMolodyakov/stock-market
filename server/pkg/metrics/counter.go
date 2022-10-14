package metrics

import "github.com/prometheus/client_golang/prometheus"

func httpResponseCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Name:      "http_response_counter",
		Help:      "Number of HTTP responses",
	}, []string{"operation", "code"})
}

func balanceActivityCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Name:      "stock_rate_counter",
		Help:      "Stock rate history",
	}, []string{"activity", "client"})
}
