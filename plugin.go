package metric

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/smallnest/rpcx/protocol"
	"net/http"
)

type PrometheusPlugin struct {
	Tag string
}

func (p PrometheusPlugin) PreHandleRequest(ctx context.Context, r *protocol.Message) error {
	std.startedCounter.WithLabelValues(r.ServicePath, r.ServiceMethod, p.Tag).Inc()
	return nil
}

func (p PrometheusPlugin) PostWriteResponse(ctx context.Context, req *protocol.Message, res *protocol.Message, err error) error {
	sp := res.ServicePath
	sm := res.ServiceMethod

	if sp == "" {
		return nil
	}

	if err != nil {
		std.handledCounter.WithLabelValues(sp, sm, p.Tag, string(Failure)).Inc()
	} else {
		std.handledCounter.WithLabelValues(sp, sm, p.Tag, string(Success)).Inc()
	}

	return nil
}

func NewPrometheusPlugin(name, tag, pattern string) *PrometheusPlugin {
	initMetrics(name)

	http.Handle(pattern, promhttp.Handler())

	return &PrometheusPlugin{
		Tag: tag,
	}
}
