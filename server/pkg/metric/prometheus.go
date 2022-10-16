package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	registry *prometheus.Registry
	handler  http.HandlerFunc
}

func NewPrometheusClient(custom bool) Prometheus {
	reg := prometheus.NewRegistry()

	if custom {
		return Prometheus{
			registry: reg,
			handler: func(w http.ResponseWriter, r *http.Request) {
				promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
			},
		}

	} else {
		return Prometheus{
			registry: reg,
			handler: func(w http.ResponseWriter, r *http.Request) {
				promhttp.Handler().ServeHTTP(w, r)
			},
		}
	}
}

func (p Prometheus) Handler() http.HandlerFunc {
	return p.handler
}

func (p Prometheus) Registry() *prometheus.Registry {
	return p.registry
}
