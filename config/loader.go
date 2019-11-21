package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
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
	if env := viper.GetString("NEON_MODE"); env != "" {
		env = strings.ToUpper(env)
		if env != NeonModeLt && env != NeonModeDev && env != NeonModeTest && env != NeonModeProduct {
			return errors.Format("NEON_MODE env should be set %s OR %s OR %s OR %s", NeonModeLt, NeonModeDev, NeonModeTest, NeonModeProduct)
		}
		var provider = viper.GetString("CONFIG_PROVIDER")
		if provider == "" {
			provider = "etcd"
		}
		var endpoint = viper.GetString("CONFIG_ENDPOINT")
		var path = viper.GetString("CONFIG_PATH")
		if path == "" {
			path = fmt.Sprintf("/config/%s/app.yml", strings.ToLower(env))
		}
		if err := LoadRemote(provider, endpoint, path); err != nil {
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

	return nil
}

func LoadRemote(provider, endpoint, path string) error {
	viper.SetConfigType("yml")
	err := viper.AddRemoteProvider(provider, endpoint, path)
	if err != nil {
		return errors.By(err)
	}

	err = viper.ReadRemoteConfig()
	if err != nil {
		return errors.By(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	return nil
}
