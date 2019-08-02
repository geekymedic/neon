package es

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivere/elastic"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/spf13/viper"
)

var (
	elasticSearch = map[string]*elastic.Client{}
)

func init() {

	type ESOptions struct {
		DSN string
	}

	neon.AddPlugin("es", func(status neon.PluginStatus, viper *viper.Viper) error {
		switch status {
		case neon.PluginLoad:

			var (
				dsnList = make(map[string]*ESOptions)
			)

			err := viper.UnmarshalKey("es", &dsnList)

			if err != nil {
				return errors.By(err)
			}

			if len(dsnList) == 0 {
				return errors.Format("es plugin used, but es config not exists.")
			}

			for name, opt := range dsnList {

				errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
				client, err := elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(opt.DSN))
				if err != nil {
					return errors.By(err)
				}
				info, code, err := client.Ping(opt.DSN).Do(context.Background())
				if err != nil {
					return errors.By(err)
				}
				fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

				esversion, err := client.ElasticsearchVersion(opt.DSN)
				if err != nil {
					return errors.By(err)
				}
				fmt.Printf("Elasticsearch version %s\n", esversion)
				elasticSearch[name] = client
			}

		}

		return nil

	})
}

func IsNotFound(err error) bool {
	if err != nil {
		e, ok := err.(*elastic.Error)
		if ok {
			if e.Status == 404 {
				return true
			}
		}
	}
	return false
}

func Use(name string) *elastic.Client {
	return elasticSearch[strings.ToLower(name)]
}
