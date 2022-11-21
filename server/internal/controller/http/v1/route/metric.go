package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrometheusClient interface {
	Handler() http.HandlerFunc
}

type prometheusRouter struct {
	prometheusClient PrometheusClient
}

func NewPrometheusRouter(client PrometheusClient) *prometheusRouter {
	return &prometheusRouter{prometheusClient: client}
}

func (pr *prometheusRouter) MetricRoute(rg *gin.RouterGroup) {
	rg.GET("/metrics", gin.WrapH(pr.prometheusClient.Handler()))
}
