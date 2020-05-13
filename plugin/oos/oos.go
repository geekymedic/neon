package oos

import (
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"

	"github.com/geekymedic/neon"
)

var ossClient = map[string]*oss.Bucket{}
var opts = map[string]struct {
	EndPoint    string
	AccessKeyId string
	KeySercet   string
	BucketName  string
}{}

func init() {
	neon.AddPlugin("oos", func(status neon.PluginStatus, viper *viper.Viper) error {
		err := viper.UnmarshalKey("oos", &opts)
		if err != nil {
			return err
		}
		for key, opt := range opts {
			key = strings.ToLower(key)
			client, err := oss.New(opt.EndPoint, opt.AccessKeyId, opt.KeySercet)
			if err != nil {
				return err
			}
			ossClient[key], err = client.Bucket(opt.BucketName)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func Use(s string) *oss.Bucket {
	return ossClient[strings.ToLower(s)]
}
