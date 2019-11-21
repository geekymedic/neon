package config

import (
	"github.com/spf13/viper"

	"github.com/geekymedic/neon/errors"
)

const (
	NeonModeLt      = "LT"
	NeonModeDev     = "DEV"
	NeonModeTest    = "TEST"
	NeonModeProduct = "PRD"
)

func init() {
	viper.BindEnv("NEON_MODE")
	viper.BindEnv("CONFIG_PROVIDER")
	viper.BindEnv("CONFIG_ENDPOINT")
	viper.BindEnv("CONFIG_PATH")
}

func Load(path *string) error {
	if env := viper.Get("NEON_MODE"); env != nil {
		switch env {
		case NeonModeLt:
			if err := LoadRemote(viper.GetString("CONFIG_PROVIDER"), viper.GetString("CONFIG_ENDPOINT"), "/config/lt/app.yml"); err != nil {
				return errors.By(err)
			}
		case NeonModeDev:
			if err := LoadRemote(viper.GetString("CONFIG_PROVIDER"), viper.GetString("CONFIG_ENDPOINT"), "/config/dev/app.yml"); err != nil {
				return errors.By(err)
			}
		case NeonModeTest:
			if err := LoadRemote(viper.GetString("CONFIG_PROVIDER"), viper.GetString("CONFIG_ENDPOINT"), "/config/test/app.yml"); err != nil {
				return errors.By(err)
			}
		case NeonModeProduct:
			if err := LoadRemote(viper.GetString("CONFIG_PROVIDER"), viper.GetString("CONFIG_ENDPOINT"), "/config/prd/app.yml"); err != nil {
				return errors.By(err)
			}
		}

		return nil
	}

	viper.AddConfigPath(*path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		return errors.By(err)
	}

	return nil
}

func LoadRemote(provider, endpoint, path string) error {
	err := viper.AddRemoteProvider(provider, endpoint, path)
	viper.SetConfigType("yml")
	if err != nil {
		return errors.By(err)
	}

	err = viper.ReadRemoteConfig()
	if err != nil {
		return errors.By(err)
	}
	return nil
}
