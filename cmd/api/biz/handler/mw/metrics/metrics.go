package metrics

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/nnieie/golanglab5/pkg/tracer"
)

// MetricsMiddleware Metrics 中间件
func MetricsMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Path())
		method := string(c.Method())

		// 继续处理请求
		c.Next(ctx)

		// 记录请求计数
		tracer.RequestCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status", c.Response.StatusCode()),
			),
		)

		// 记录请求延迟
		duration := time.Since(start).Seconds()
		tracer.RequestDuration.Record(ctx, duration,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
			),
		)

		// 记录错误
		if c.Response.StatusCode() >= consts.StatusBadRequest {
			tracer.ErrorCounter.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("method", method),
					attribute.String("path", path),
					attribute.Int("status", c.Response.StatusCode()),
				),
			)
		}
	}
}
