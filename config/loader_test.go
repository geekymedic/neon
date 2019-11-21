package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/stretchr/testify/require"
)

func TestLoadRemote(t *testing.T) {
	viper.BindEnv("NEON_MODE")
	viper.BindEnv("CONFIG_PROVIDER")
	viper.BindEnv("CONFIG_ENDPOINT")
	viper.BindEnv("CONFIG_PATH")

	t.Run("it is ok", func(t *testing.T) {
		os.Setenv("NEON_MODE", "DEV")
		// os.Setenv("CONFIG_PROVIDER", "etcd")
		os.Setenv("CONFIG_ENDPOINT", "http://192.168.0.202:12379")
		// os.Setenv("CONFIG_PATH", "/config/dev/app.yml")
		err := Load(nil)
		require.Nil(t, err)
		viper.Get("servers.etcd-system-bff-demo")
		t.Log(viper.AllKeys())
	})

	t.Run("it is fail", func(t *testing.T) {
		os.Setenv("NEON_MODE", "XXXX")
		err := Load(nil)
		require.NotNil(t, err)
	})
}
