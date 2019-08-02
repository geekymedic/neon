package config

import (
	"github.com/geekymedic/neon/errors"
	"github.com/spf13/viper"
)

func Load(path *string) error {
	viper.AddConfigPath(*path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()

	if err != nil {
		return errors.By(err)
	}

	return nil
}
