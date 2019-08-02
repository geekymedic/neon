package bff

import (
	"flag"
	"fmt"
	"github.com/geekymedic/neon/version"
	"net"
	"net/http"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/config"
	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"

	"github.com/spf13/viper"
)

func Main() error {
	fmt.Println(version.Version())

	path := flag.String("c", "config", "config path")
	flag.Parse()
	err := config.Load(path)

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
		logger.Errorf("listen %s fail, %v\n", address, err)
		return errors.By(err)
	} else {
		logger.Infof("listen %s", l.Addr())
	}

	return http.Serve(l, _engine)
}
