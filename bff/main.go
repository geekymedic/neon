package bff

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/config"
	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/version"

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

	srv := &http.Server{
		Addr:    address,
		Handler: _engine,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		logger.With("address", srv.Addr).Info("Start listen...")
		err = srv.ListenAndServe()
		if err != nil {
			logger.Errorf("listen %s fail, %v\n", address, err)
			err = errors.By(err)
		} else {
			logger.Infof("listen %s", srv.Addr)
		}
		c<- syscall.SIGTERM
	}()
	<-c

	logger.Info("Shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	return errors.Wrap(err)
}
