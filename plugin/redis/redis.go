package redis

import (
	"strings"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/spf13/viper"

	"github.com/go-redis/redis"
)

var (
	redisList = map[string]*redis.Client{}
)

func init() {

	type RedisOptions struct {
		DSN      string
		Password string
		DB       int
	}

	neon.AddPlugin("redis", func(status neon.PluginStatus, viper *viper.Viper) error {
		switch status {
		case neon.PluginLoad:

			var (
				dsnList = make(map[string]*RedisOptions)
			)

			err := viper.UnmarshalKey("redis", &dsnList)

			if err != nil {
				return errors.By(err)
			}

			if len(dsnList) == 0 {
				return errors.Format("redis plugin used, but redis config not exists.")
			}

			for name, opt := range dsnList {

				client := redis.NewClient(&redis.Options{
					Addr:     opt.DSN,
					Password: opt.Password,
					DB:       opt.DB,
				})

				redisList[strings.ToLower(name)] = client
			}

		}

		return nil

	})
}

func Use(name string) *redis.Client {
	return redisList[strings.ToLower(name)]
}
