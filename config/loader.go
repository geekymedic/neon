package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
)

const (
	NeonModeLt      = "LT"
	NeonModeDev     = "DEV"
	NeonModeTest    = "TEST"
	NeonModeProduct = "PRD"
)

const (
	NeonMode           = "NEON_MODE"
	NeonConfigProvider = "NEON_CONFIG_PROVIDER"
	NeonConfigEndpoint = "NEON_CONFIG_ENDPOINT"
	NeonConfigPath     = "NEON_CONFIG_PATH"
	NeonConfigSecret   = "NEON_CONFIG_SECRET"
)

func init() {
	viper.BindEnv(NeonMode)
	viper.BindEnv(NeonConfigProvider)
	viper.BindEnv(NeonConfigEndpoint)
	viper.BindEnv(NeonConfigPath)
	viper.BindEnv(NeonConfigSecret)
}

var RemoterViper = viper.New()
var LocalViper = viper.New()

// Must compatible old version, don't' ignore the function
func Load(path *string) error {
	if env := viper.GetString(NeonMode); env != "" {
		viper.SupportedRemoteProviders = []string{"etcd", "apollo"}
		env = strings.ToUpper(env)
		if env != NeonModeLt && env != NeonModeDev && env != NeonModeTest && env != NeonModeProduct {
			return errors.Format("NEON_MODE env should be set %s OR %s OR %s OR %s", NeonModeLt, NeonModeDev, NeonModeTest, NeonModeProduct)
		}
		if err := LoadRemote(*path); err != nil {
			return errors.By(err)
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
	LocalViper = viper.GetViper()
	viper.WatchConfig()
	return nil
}

func LoadRemote(localPath string) error {
	// read the local config
	{
		LocalViper.AddConfigPath(localPath)
		LocalViper.SetConfigName("config")
		LocalViper.SetConfigType("yml")
		if err := LocalViper.ReadInConfig(); err != nil {
			return errors.Wrap(err)
		}
		if err := LocalViper.ReadInConfig(); err != nil {
			return errors.Wrap(err)
		}
		LocalViper.WatchConfig()
	}

	viper.SetConfigType("yml")
	logger.With(NeonConfigProvider, viper.Get(NeonConfigProvider),
		NeonConfigEndpoint, viper.Get(NeonConfigEndpoint),
		NeonConfigPath, viper.Get(NeonConfigPath),
		NeonConfigSecret, viper.Get(NeonConfigSecret)).Debug("config env trace")
	err := viper.AddSecureRemoteProvider(viper.GetString(NeonConfigProvider),
		viper.GetString(NeonConfigEndpoint),
		viper.GetString(NeonConfigPath),
		viper.GetString(NeonConfigSecret))
	if err != nil {
		return errors.Wrap(err)
	}
	if err = viper.ReadRemoteConfig(); err != nil {
		return errors.Wrap(err)
	}
	if err = viper.MergeConfigMap(LocalViper.AllSettings()); err != nil {
		return errors.Wrap(err)
	}
	if err = viper.GetViper().WatchRemoteConfigOnChannel(); err != nil {
		return errors.By(err)
	}
	// Bind remote viper
	RemoterViper = viper.GetViper()
	fmt.Println("Pass load remote config from ", NeonConfigProvider)
	return nil
}
