package metric

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	HTTPResponseCounter       *prometheus.CounterVec
	ResponseDurationHistogram *prometheus.HistogramVec
}

func NewMetric(registry *prometheus.Registry) Metric {
	m := &Metric{}

	m.HTTPResponseCounter = httpResponseCounter()
	registry.MustRegister(m.HTTPResponseCounter)

	m.ResponseDurationHistogram = responseDurationHistogram()
	registry.MustRegister(m.ResponseDurationHistogram)

	return *m
}
