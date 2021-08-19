package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	startedCounter   *prom.CounterVec
	handledCounter   *prom.CounterVec
	handledHistogram *prom.HistogramVec
}

var std *Metrics

func initMetrics(name string) {
	std = &Metrics{
		startedCounter: prom.NewCounterVec(
			prom.CounterOpts{
				Name: name+"_started_total",
				Help: "Total number started on the server.",
			}, []string{"service", "method", "tag"}),
		handledCounter: prom.NewCounterVec(
			prom.CounterOpts{
				Name: name+"_handled_total",
				Help: "Total number completed on the server, regardless of success or failure.",
			}, []string{"service", "method", "tag", "code"}),
		handledHistogram: prom.NewHistogramVec(
			prom.HistogramOpts{
				Name:    name+"_handling_seconds",
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
