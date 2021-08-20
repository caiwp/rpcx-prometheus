package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	startedCounter   *prom.CounterVec
	handledCounter   *prom.CounterVec
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
	}

	prom.MustRegister(std.startedCounter, std.handledCounter)
}

type Code string

const (
	Success Code = "Success"
	Failure Code = "Failure"
)
