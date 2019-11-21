package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/stretchr/testify/assert"
)

func TestLoadRemote(t *testing.T) {
	os.Setenv("NEON_MODE", "DEV")
	os.Setenv("CONFIG_PROVIDER", "etcd")
	os.Setenv("CONFIG_ENDPOINT", "http://192.168.0.202:12379")
	os.Setenv("CONFIG_PATH", "/config/dev/app.yml")

	viper.BindEnv("NEON_MODE")
	viper.BindEnv("CONFIG_PROVIDER")
	viper.BindEnv("CONFIG_ENDPOINT")
	viper.BindEnv("CONFIG_PATH")
	err := Load(nil)
	t.Log(err)
	assert.Nil(t, err)
	viper.Get("servers.etcd-system-bff-demo")
	t.Log(viper.AllKeys())
}
