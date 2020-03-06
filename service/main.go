package service

import (
	"fmt"
	"github.com/geekymedic/neon/service/middleware"
	"net"
	"sync"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/config"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var _rpcServer *grpc.Server

func Main() error {
	fmt.Println(version.Version())

	var path string
	cmd := cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {

			err := config.Load(&path)

			if err != nil {
				return err
			}

			var (
				address = viper.GetString("Address")
			)

			err = neon.LoadPlugins(viper.GetViper())

			if err != nil {
				return err
			}
			for _, opt := range beforeAppRun {
				if err := opt(); err != nil {
					logger.With("err", err).Error("fail to execute before app run hook")
					return err
				}
			}

			defer func() {
				for _, opt := range beforeAppExit {
					if err := opt(); err != nil {
						logger.With("err", err).Error("fail to execute before app exit hook")
					}
				}
			}()
			l, err := net.Listen("tcp", address)

			if err != nil {
				logger.Errorf("listen %s fail,%s", address, err)
				return err
			} else {
				logger.Info("listen %s", l.Addr())
			}
			Server()
			err = _rpcServer.Serve(l)

			if err != nil {
				logger.Errorf("%s", err)
				return err
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&_flags.Debug, "debug", false, "use debug mode")
	flags.StringVarP(&path, "config", "c", "config", "config path")
	return cmd.Execute()

}

var once sync.Once

func Server() *grpc.Server {
	once.Do(func() {
		_rpcServer = grpc.NewServer(grpc.UnaryInterceptor(middleware.GrpcServerChainMiddleware()))
	})
	return _rpcServer
}
