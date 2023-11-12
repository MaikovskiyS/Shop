package metrics

import (
	"myproject/internal/server/router"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RequestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "shop",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"status"})

func ObserveRequest(d time.Duration, status int) {
	RequestMetrics.WithLabelValues(strconv.Itoa(status)).Observe(d.Seconds())
}
func Register(r *router.Router) {
	r.Handle("/metrics", promhttp.Handler())

}
