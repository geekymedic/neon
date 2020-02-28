package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	viper.BindEnv(NeonMode)
	viper.BindEnv(NeonConfigProvider)
	viper.BindEnv(NeonConfigEndpoint)
	viper.BindEnv(NeonCo	nfigPath)
	viper.BindEnv(NeonConfigSecret)
	os.Setenv(NeonMode, NeonModeDev)
	os.Setenv(NeonConfigProvider, "etcd")
	os.Setenv(NeonConfigEndpoint, "http://192.168.0.202:12379,http://192.168.0.202:12379")
	os.Setenv(NeonConfigPath, "/config/dev")
	os.Setenv(NeonConfigSecret, "root@123456")
	var path = "."
	require.Nil(t, Load(&path))

	t.Run("etcd", func(t *testing.T) {
		for {
			err := viper.GetViper().WatchRemoteConfigOnChannel()
			for key, value := range viper.AllSettings() {
				fmt.Println(key, value)
			}
			require.Nil(t, err)
			time.Sleep(time.Second * 3)
			fmt.Println(viper.GetString("Metrics.Address"))
			fmt.Println("--------------------------------")
		}
	})

	// t.Run("apollo", func(t *testing.T) {
	// 	viper.SetConfigFile("yml")
	// 	viper.SupportedRemoteProviders = []string{"etcd", "apollo"}
	// 	var args = []struct {
	// 		provider string
	// 		endpoint string
	// 		path     string
	// 		secret   string
	// 	}{
	// 		{
	// 			provider: "apollo",
	// 			endpoint: "localhost:8080",
	// 			path:     "DEV.2.yml",
	// 			secret:   "SampleApp",
	// 		},
	// 	}
	// 	for _, arg := range args {
	// 		err := viper.AddSecureRemoteProvider(arg.provider, arg.endpoint, arg.path, arg.secret)
	// 		require.Nil(t, err)
	// 	}
	// 	err := viper.ReadRemoteConfig()
	// 	require.Nil(t, err)
	// 	t.Log(viper.AllKeys())
	// 	err = viper.GetViper().WatchRemoteConfigOnChannel()
	// 	require.Nil(t, err)
	// })
}
