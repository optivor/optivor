package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "optivor_requests_total",
			Help: "Total number of HTTP requests served by Optivor.",
		},
		[]string{"status", "format", "fit"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "optivor_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	cacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "optivor_cache_hits_total",
			Help: "Total number of cache hits.",
		},
	)

	cacheMissesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "optivor_cache_misses_total",
			Help: "Total number of cache misses.",
		},
	)

	transformDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "optivor_transform_duration_seconds",
			Help:    "Pipeline image transformation duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"format", "fit"},
	)
)

func init() {
	prometheus.MustRegister(
		requestsTotal,
		requestDuration,
		cacheHitsTotal,
		cacheMissesTotal,
		transformDuration,
	)
}

// MetricsHandler returns an http.Handler for Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
