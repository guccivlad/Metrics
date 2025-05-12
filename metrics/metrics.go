package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	*prometheus.Registry
	LogCalls    prometheus.Counter
	SuccesLogs  prometheus.Counter
	FailedLogs  prometheus.Counter
	ReqDuration prometheus.Histogram
}

func NewMetrics() (*Metrics, error) {
	m := &Metrics{
		Registry: prometheus.NewRegistry(),
	}

	m.LogCalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "log_calls_total",
		Help: "total number of /log calls",
	})

	m.SuccesLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "succes_log_calls",
		Help: "succes /log calls",
	})

	m.FailedLogs = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_log_calls",
		Help: "failed /log calls",
	})

	m.ReqDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "log_request_duration_seconds",
		Help:    "Duration of /log handler in seconds",
		Buckets: prometheus.DefBuckets,
	})

	for _, metric := range []prometheus.Collector{
		m.LogCalls, m.FailedLogs, m.SuccesLogs, m.ReqDuration,
	} {
		if err := m.Registry.Register(metric); err != nil {
			return nil, err
		}
	}

	return m, nil
}
