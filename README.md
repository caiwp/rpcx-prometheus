# Go-rpcx-prometheus
go rpcx prometheus plugin

```
rpcx + prometheus 监控指标
主要提供 http 服务 /metrics 路由让 prometheus 采集

import (
    prom "github.com/caiwp/rpcx-prometheus"
)

s.Plugins.Add(prom.NewPrometheusPlugin("service", "demoTag", "/metrics"))

```