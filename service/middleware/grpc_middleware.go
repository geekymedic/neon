package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/logger/extend"
	"github.com/geekymedic/neon/metrics/prometheus"
	"github.com/geekymedic/neon/version"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func grpcLogMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var (
			start               = time.Now()
			log   logger.Logger = extend.DefGrpLog()
			ses                 = neon.CreateSessionFromGrpcIncomingContext(ctx)
		)

		log = log.With(sessionTraceLog(ses)...).
			With("pro_name", version.PRONAME,
				"service", fmt.Sprintf("%T", info.Server), "method", info.FullMethod)
		buf, _ := json.Marshal(req)
		log = log.With("inbound", string(buf))

		resp, err = handler(ctx, req)

		log = log.With("latency", fmt.Sprintf("%v", time.Now().Sub(start)))
		if err != nil {
			log.With("err", err).Error("grpc request trace")
			err = errors.WithMessage(err, "pro_name:%s, service:%s, method:%s", version.PRONAME, fmt.Sprintf("%T", info.Server), info.FullMethod)
		} else {
			buf, _ = json.Marshal(resp)
			log.With("outbound", string(buf)).Info("grpc request trace")
		}
		return resp, err
	}
}

func grpcMetricsMiddleware() grpc.UnaryServerInterceptor {
	if os.Getenv("NEON_MODE") == "test" {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			resp, err = handler(ctx, req)
			return resp, err
		}
	}
	var qpsMetrics = prometheus.MustCounterWithLabelNames("request_qps", "service", "method", "ret")
	var latency = prometheus.MustSummaryWithLabelNames("request_latency", map[float64]float64{
		0.5: 0.005, 0.9: 0.01, 0.99: 0.001},
		"service", "method")
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var (
			start = time.Now()
		)

		resp, err = handler(ctx, req)

		// qps
		{
			if err == nil {
				qpsMetrics.With(fmt.Sprintf("%T", info.Server), info.FullMethod, "ok").Inc()
			} else {
				qpsMetrics.With(fmt.Sprintf("%T", info.Server), info.FullMethod, "fail").Inc()
			}
		}

		// latency
		{
			elasp := float64(time.Since(start)) / float64(time.Millisecond)
			latency.With(fmt.Sprintf("%T", info.Server), info.FullMethod).Observe(elasp)
		}

		return resp, err
	}
}

func GrpcServerChainMiddleware(serverInterceptor ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	chain := []grpc.UnaryServerInterceptor{grpcMetricsMiddleware(), grpcLogMiddleware()}
	chain = append(chain, serverInterceptor...)
	return grpc_middleware.ChainUnaryServer(chain...)
}

func sessionTraceLog(ses *neon.Session) []interface{} {
	return []interface{}{
		"_uid", ses.Uid,
		"_token", ses.Token,
		"_trace", ses.Trace,
		"_sequence", ses.Sequence,
		"_time", ses.Time,
		"_storeId", ses.StoreId,
		"_clientIp", ses.ClientIp,
	}
}
