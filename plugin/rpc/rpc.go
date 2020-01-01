package rpc

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/logger"
)

var connections = map[string]*grpc.ClientConn{}
var lock = &sync.RWMutex{}

func init() {

	neon.AddPlugin("rpc_server", func(status neon.PluginStatus, viper *viper.Viper) error {
		switch status {
		case neon.PluginLoad:
			// servers := viper.GetStringMapString("servers")
			// for name, address := range servers {
			// 	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpcClientLog()))
			// 	if err != nil {
			// 		return errors.By(err)
			// 	}
			// 	connections[name] = conn
			// }
		}
		return nil

	})
}

func GetConnection(name string) *grpc.ClientConn {
	lock.RLock()
	conn, ok := connections[name]
	lock.RUnlock()
	if ok {
		return conn
	}
	server := fmt.Sprintf("%s", name)
	address := viper.GetString(server)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpcClientLog()))
	if err != nil {
		logger.With("address", address).Error(err)
		return nil
	}
	lock.Lock()
	connections[strings.ToLower(name)] = conn
	defer lock.Unlock()
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
