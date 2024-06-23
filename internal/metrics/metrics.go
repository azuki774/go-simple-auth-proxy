package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	ExporterPort string
}

var (
	AccessCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "go_simple_auth_proxy",
		Name:      "access_count",
		Help:      "access_count for auth response",
	},
		[]string{"response"},
	)
)

func NewMetrics(port string) *MetricsServer {
	return &MetricsServer{ExporterPort: port}
}

func (m *MetricsServer) Start() (err error) {
	prometheus.MustRegister(AccessCounterVec)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%s", m.ExporterPort), promhttp.Handler())
}
