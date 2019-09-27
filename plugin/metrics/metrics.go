package metrics

import (
	_ "net/http/pprof"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/metrics/prometheus"
	"github.com/spf13/viper"
)

func init() {

	neon.AddPlugin("metrics", func(status neon.PluginStatus, viper *viper.Viper) error {
		switch status {
		case neon.PluginLoad:

			addr := viper.GetString("Metrics.Address")
			if addr == "" {
				logger.Warn("not found Metrics.Address config")
				return nil
			}
			path := viper.GetString("Metrics.Path")
			if path == "" {
				path = "/metrics"
			}

			if len(addr) == 0 {
				return errors.Format("load Metrics.Address fail, empty address.")
			}
			go func() {
				logger.With("lis", addr, "path", path).Info("start metrics server")
				err := prometheus.StartMetricsServer(addr, path)
				if err != nil {
					logger.Error("fail to start metrics server", "err", err)
				}
			}()
		}

		return nil
	})

}
