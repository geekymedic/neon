package rpc

import (
	"context"
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc/codes"
	"path"
	"strings"
	"time"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var connections = map[string]*grpc.ClientConn{}

func init() {

	neon.AddPlugin("rpc_server", func(status neon.PluginStatus, viper *viper.Viper) error {
		switch status {
		case neon.PluginLoad:
			servers := viper.GetStringMapString("servers")
			for name, address := range servers {
				conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryMiddle()...)),
					grpc.WithUnaryInterceptor(grpcClientLog()))
				if err != nil {
					return errors.By(err)
				}
				connections[name] = conn
			}
		}
		return nil

	})
}

func GetConnection(name string) *grpc.ClientConn {
	return connections[strings.ToLower(name)]
}

func MockGrpcClientLog() grpc.UnaryClientInterceptor {
	return grpcClientLog()
}

func grpcClientLog() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(method)[1:]
		serviceMethod := path.Base(method)
		ses := neon.CreateSessionFromGrpcOutgoingContext(ctx)
		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		log := logger.With(sessionTraceLog(ses)...).With("grpc.service", service, "grpc.method", serviceMethod, "latency", fmt.Sprintf("%v", time.Now().Sub(startTime)))
		if err != nil {
			log.With("err", err).Error("finished client unary call")
		} else {
			log.Info("finished client unary call")
		}
		return err
	}
}

func retryMiddle() []grpc_retry.CallOption {
	retryOpt := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Unavailable),
	}
	return retryOpt
}

func sessionTraceLog(ses *neon.Session) []interface{} {
	return []interface{}{
		"_uid", ses.Uid,
		"_token", ses.Token,
		"_trace", ses.Trace,
		"_sequence", ses.Sequence,
		"_time", ses.Time,
		"_storeId", ses.StoreId,
	}
}
