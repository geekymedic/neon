package service_mock

import (
	"os"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/service"
	"github.com/spf13/viper"
)

func InitEnv(configName, configDir string) error {
	err := os.Setenv("NEON_MODE", "test")
	if err != nil {
		return err
	}
	logger.SetLevel(logger.DebugLevel)

	viper.SetConfigName(configName)
	viper.AddConfigPath(configDir)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = neon.LoadPlugins(viper.GetViper())
	if err != nil {
		return err
	}

	for _, opt := range service.GetBeforeAppRun() {
		if err = opt(); err != nil {
			return err
		}
	}

	return nil
}
