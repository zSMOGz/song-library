package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"song-library/internal/constants"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: constants.MetricHTTPRequestsTotal,
			Help: constants.MetricHTTPRequestsHelp,
		},
		[]string{
			constants.MetricLabelMethod,
			constants.MetricLabelEndpoint,
			constants.MetricLabelStatus,
		},
	)
)
