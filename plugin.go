package metric

import (
	"context"
	"github.com/smallnest/rpcx/server"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/smallnest/rpcx/protocol"
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

	t := ctx.Value(server.StartRequestContextKey).(int64)
	if t > 0 {
		t = time.Now().UnixNano()-t
		if t < 10*time.Minute.Nanoseconds() {
			std.handledHistogram.WithLabelValues(sp, sm, p.Tag).Observe(float64(t))
		}
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
