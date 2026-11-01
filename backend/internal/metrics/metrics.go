package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds the registered Prometheus metric collectors for FlagManagment.
type Metrics struct {
	Registry *prometheus.Registry

	EvaluationsTotal      *prometheus.CounterVec
	EvaluationDuration    *prometheus.HistogramVec
	CacheOperationsTotal  *prometheus.CounterVec
	ActiveFlags           *prometheus.GaugeVec
	StaleFlags            *prometheus.GaugeVec
	ChangeRequestsTotal   *prometheus.CounterVec
	HttpRequestsTotal     *prometheus.CounterVec
	HttpRequestDuration   *prometheus.HistogramVec
}

var globalMetrics *Metrics

// Init initializes the Prometheus metrics registry and instruments collectors.
func Init() *Metrics {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		Registry: reg,

		EvaluationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "flagmanagment",
				Subsystem: "evaluator",
				Name:      "evaluations_total",
				Help:      "Total number of feature flag evaluations performed.",
			},
			[]string{"project_id", "environment", "flag_key", "variant", "reason"},
		),

		EvaluationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "flagmanagment",
				Subsystem: "evaluator",
				Name:      "evaluation_duration_seconds",
				Help:      "Latency of feature flag evaluations in seconds.",
				Buckets:   []float64{0.000005, 0.00001, 0.00005, 0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05},
			},
			[]string{"project_id", "environment"},
		),

		CacheOperationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "flagmanagment",
				Subsystem: "cache",
				Name:      "operations_total",
				Help:      "Total number of cache operations by type (hit, miss, set, evict).",
			},
			[]string{"operation"},
		),

		ActiveFlags: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "flagmanagment",
				Subsystem: "flags",
				Name:      "active_total",
				Help:      "Current count of active feature flags by project.",
			},
			[]string{"project_id"},
		),

		StaleFlags: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "flagmanagment",
				Subsystem: "flags",
				Name:      "stale_total",
				Help:      "Current count of stale feature flags flagged by governance bot.",
			},
			[]string{"project_id"},
		),

		ChangeRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "flagmanagment",
				Subsystem: "governance",
				Name:      "change_requests_total",
				Help:      "Total number of change requests processed by status.",
			},
			[]string{"project_id", "status"},
		),

		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "flagmanagment",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests processed by endpoint and status.",
			},
			[]string{"method", "path", "status"},
		),

		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "flagmanagment",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "HTTP request latency distribution in seconds.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}

	reg.MustRegister(
		m.EvaluationsTotal,
		m.EvaluationDuration,
		m.CacheOperationsTotal,
		m.ActiveFlags,
		m.StaleFlags,
		m.ChangeRequestsTotal,
		m.HttpRequestsTotal,
		m.HttpRequestDuration,
	)

	globalMetrics = m
	return m
}

// Global returns the initialized global metrics instance, initializing if necessary.
func Global() *Metrics {
	if globalMetrics == nil {
		globalMetrics = Init()
	}
	return globalMetrics
}

// Handler returns an http.Handler that exposes Prometheus metrics.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
}

// RecordEvaluation records a flag evaluation occurrence.
func (m *Metrics) RecordEvaluation(projectID, env, flagKey, variant, reason string, duration time.Duration) {
	if m == nil {
		return
	}
	m.EvaluationsTotal.WithLabelValues(projectID, env, flagKey, variant, reason).Inc()
	m.EvaluationDuration.WithLabelValues(projectID, env).Observe(duration.Seconds())
}

// RecordCacheHit records a cache hit.
func (m *Metrics) RecordCacheHit() {
	if m == nil {
		return
	}
	m.CacheOperationsTotal.WithLabelValues("hit").Inc()
}

// RecordCacheMiss records a cache miss.
func (m *Metrics) RecordCacheMiss() {
	if m == nil {
		return
	}
	m.CacheOperationsTotal.WithLabelValues("miss").Inc()
}
