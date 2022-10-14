package metrics

import "github.com/prometheus/client_golang/prometheus"

func responseDurationHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "client",
		Name:      "balance_response_duration_histogram",
		Help:      "Balance response duration (ms)",
		Buckets:   []float64{10, 50, 90, 130, 170, 210, 250, 290, 330},
		// This is same as prometheus.LinearBuckets(10, 40, 9)
		// 9 buckets starting from 10 increased by 40
	}, []string{"operation"})
}
