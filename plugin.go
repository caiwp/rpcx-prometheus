package metric

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
)

const key = "prom_unit"

type PrometheusPlugin struct {
	Tag string
}

func (p PrometheusPlugin) PreHandleRequest(ctx context.Context, r *protocol.Message) error {
	u := NewUnit(r.ServicePath, r.ServiceMethod, p.Tag)
	c, ok := ctx.(*share.Context)
	if ok {
		c.SetValue(key, u)
	}

	return nil
}

func (p PrometheusPlugin) PostWriteResponse(ctx context.Context, req *protocol.Message, res *protocol.Message, err error) error {
	if c, ok := ctx.(*share.Context); ok {
		u, ok := c.Value(key).(*Unit)
		if ok {
			if err != nil {
				u.Handle(req.ServicePath, req.ServiceMethod, p.Tag, Failure)
			} else {
				u.Handle(req.ServicePath, req.ServiceMethod, p.Tag, Success)
			}
		}
	}

	return nil
}

func NewPrometheusPlugin(tag, pattern string) *PrometheusPlugin {
	initMetrics()

	http.Handle(pattern, promhttp.Handler())

	return &PrometheusPlugin{
		Tag: tag,
	}
}
