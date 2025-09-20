package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics содержит все метрики приложения
type Metrics struct {
	Registry prometheus.Registerer
	Gatherer prometheus.Gatherer

	// HTTP
	HTTPRequestCount    *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPInflight        prometheus.Gauge

	// DB
	DBQueryDuration *prometheus.HistogramVec
	DBQueryErrors   *prometheus.CounterVec

	// Kafka
	KafkaMessagesConsumed *prometheus.CounterVec
	KafkaMessagesDLQ      *prometheus.CounterVec
	KafkaConsumerRetries  *prometheus.CounterVec

	// Redis / Cache
	CacheHits   prometheus.Counter
	CacheMisses prometheus.Counter
}

// New создает метрики и регистрирует их в регистри
func New(reg prometheus.Registerer, gatherer prometheus.Gatherer) *Metrics {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}
	if gatherer == nil {
		gatherer = prometheus.DefaultGatherer
	}

	m := &Metrics{
		Registry: reg,
		Gatherer: gatherer,

		HTTPRequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "HTTP requests total",
			},
			[]string{"handler", "method", "code", "service"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request durations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"handler", "method", "service"},
		),
		HTTPInflight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "http_inflight_requests",
			Help: "Number of in-flight HTTP requests",
		}),
		DBQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "DB query duration seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
			},
			[]string{"query", "service"},
		),
		DBQueryErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_query_errors_total",
				Help: "DB query errors total",
			},
			[]string{"query", "service"},
		),
		KafkaMessagesConsumed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kafka_messages_consumed_total",
				Help: "Kafka messages consumed",
			},
			[]string{"topic", "outcome"},
		),
		KafkaMessagesDLQ: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kafka_messages_dlq_total",
				Help: "Kafka messages sent to DLQ",
			},
			[]string{"topic"},
		),
		KafkaConsumerRetries: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kafka_consumer_retries_total",
				Help: "Kafka consumer retries",
			},
			[]string{"topic"},
		),
		CacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total cache hits",
		}),
		CacheMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total cache misses",
		}),
	}

	// регистрация
	reg.MustRegister(
		m.HTTPRequestCount,
		m.HTTPRequestDuration,
		m.HTTPInflight,
		m.DBQueryDuration,
		m.DBQueryErrors,
		m.KafkaMessagesConsumed,
		m.KafkaMessagesDLQ,
		m.KafkaConsumerRetries,
		m.CacheHits,
		m.CacheMisses,
	)

	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.Gatherer, promhttp.HandlerOpts{})
}
