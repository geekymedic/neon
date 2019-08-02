package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/geekymedic/neon/bff/types"
	"github.com/geekymedic/neon/metrics/prometheus"
	"github.com/gin-gonic/gin"
)

func MetricsMiddleWare() func(ctx *gin.Context) {
	if os.Getenv("neon_MODE") == "test" {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}
	var qpsMetrics = prometheus.MustCounterWithLabelNames("request_qps", "method", "host", "path", "status")
	var latencyCounterMetrics = prometheus.MustGagueWithLabelNames("request_gague_latency", "method", "host", "path", "status")
	var latencyMetrics = prometheus.MustSummaryWithLabelNames("request_latency", map[float64]float64{
		0.5: 0.005, 0.9: 0.01, 0.99: 0.001},
		"method", "path", "status")
	var responseMetrics = prometheus.MustCounterWithLabelNames("response_status_code", "code", "path", "size")
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		elasp := float64(time.Since(start)) / float64(time.Millisecond)

		method := ctx.Request.Method
		host := ctx.Request.Host
		path := ctx.Request.URL.Path
		status := fmt.Sprintf("%d", ctx.Writer.Status())

		// qps
		{
			qpsMetrics.With(method, host, path, status).Inc()
		}

		// latency
		{
			latencyMetrics.With(method, path, status).Observe(elasp)
			latencyCounterMetrics.With(method, host, path, status).Set(elasp)
		}

		// response code
		{
			value, ok := ctx.Get(types.ResponseStatusCode)
			if !ok {
				return
			}
			code, ok := value.(int)
			if !ok {
				return
			}
			responseMetrics.With(fmt.Sprintf("%d", code), ctx.Request.URL.Path, fmt.Sprintf("%d", ctx.Writer.Size())).Inc()
		}
	}
}
