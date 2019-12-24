package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	viper.BindEnv("NEON_MODE")
	viper.BindEnv("CONFIG_PROVIDER")
	viper.BindEnv("CONFIG_ENDPOINT")
	viper.BindEnv("CONFIG_PATH")

	t.Run("it is ok", func(t *testing.T) {
		viper.SetConfigType("yml")
		var args = []struct {
			provider string
			endpoint string
			path     string
		}{
			{
				provider: "etcd",
				endpoint: "http://127.0.0.1:2379",
				path:     "/config",
			},
		}
		for _, arg := range args {
			err := viper.AddRemoteProvider(arg.provider, arg.endpoint, arg.path)
			require.Nil(t, err)
		}

		err := viper.ReadRemoteConfig()
		require.Nil(t, err)
		t.Log(viper.AllKeys())
		err = viper.GetViper().WatchRemoteConfigOnChannel()
		require.Nil(t, err)
	})
}
