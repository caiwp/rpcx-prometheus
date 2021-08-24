package metric

import (
	"context"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"time"
)

const (
	Success = "Success"
	Failure = "Failure"
)

type Plugin struct {
	tag string
}

func (p *Plugin) PostReadRequest(ctx context.Context, r *protocol.Message, e error) error {
	sp := r.ServicePath
	sm := r.ServiceMethod

	if sp == "" {
		return nil
	}

	std.startedCounter.WithLabelValues(sp, sm, p.tag).Inc()

	return nil
}

func (p *Plugin) PostWriteResponse(ctx context.Context, req *protocol.Message, res *protocol.Message, e error) error {
	sp := res.ServicePath
	sm := res.ServiceMethod

	if sp == "" {
		return nil
	}

	if e != nil {
		std.handledCounter.WithLabelValues(sp, sm, p.tag, Failure).Inc()
	} else {
		std.handledCounter.WithLabelValues(sp, sm, p.tag, Success).Inc()
	}

	t := ctx.Value(server.StartRequestContextKey).(int64)
	if t > 0 {
		t = time.Now().UnixNano() - t
		if t < 10*time.Minute.Nanoseconds() { // impossible 10 minute
			std.handledHistogram.WithLabelValues(sp, sm, p.tag).Observe(time.Duration(t).Seconds())
		}
	}

	return nil
}

func NewPlugin(tag string) *Plugin {
	return &Plugin{tag: tag}
}
