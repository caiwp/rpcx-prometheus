package metric

import (
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	startedCounter   *prom.CounterVec
	handledCounter   *prom.CounterVec
	handledHistogram *prom.HistogramVec
}

var std *Metrics

func initMetrics() {
	std = &Metrics{
		startedCounter: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "server_started_total",
				Help: "Total number started on the server.",
			}, []string{"service", "method", "tag"}),
		handledCounter: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "server_handled_total",
				Help: "Total number completed on the server, regardless of success or failure.",
			}, []string{"service", "method", "tag", "code"}),
		handledHistogram: prom.NewHistogramVec(
			prom.HistogramOpts{
				Name:    "server_handling_seconds",
				Help:    "Histogram of response latency (seconds) that had been application-level handled by the server.",
				Buckets: prom.DefBuckets,
			}, []string{"service", "method", "tag"}),
	}

	prom.MustRegister(std.startedCounter, std.handledCounter, std.handledHistogram)
}

type Code string

const (
	Success Code = "Success"
	Failure Code = "Failure"
)

type Unit struct {
	startAt time.Time
}

func NewUnit(service, method, tag string) *Unit {
	std.startedCounter.WithLabelValues(service, method, tag).Inc()
	return &Unit{startAt: time.Now()}
}

func (u *Unit) Handle(service, method, tag string, code Code) {
	std.handledCounter.WithLabelValues(service, method, tag, string(code)).Inc()
	std.handledHistogram.WithLabelValues(service, method, tag).Observe(time.Since(u.startAt).Seconds())
}
