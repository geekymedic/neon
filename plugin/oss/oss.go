package oss

import (
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"

	"github.com/geekymedic/neon"
)

var ossClient = map[string]*oss.Client{}
var opts = map[string]struct {
	EndPoint    string
	AccessKeyId string
	KeySercet   string
	BucketName  string
}{}

func init() {
	neon.AddPlugin("oss", func(status neon.PluginStatus, viper *viper.Viper) error {
		err := viper.UnmarshalKey("oss", &opts)
		if err != nil {
			return err
		}
		for key, opt := range opts {
			key = strings.ToLower(key)
			client, err := oss.New(opt.EndPoint, opt.AccessKeyId, opt.KeySercet)
			if err != nil {
				return err
			}
			ossClient[key] = client
		}
		return nil
	})
}

func Use(s string) *oss.Client {
	return ossClient[strings.ToLower(s)]
}
