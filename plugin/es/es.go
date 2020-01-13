package es

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivere/elastic"

	"crypto/tls"
	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"time"
)

var (
	elasticSearch = map[string]*elastic.Client{}
)

func init() {

	type ESOptions struct {
		DSN string
		USER string
		PASSWORD string
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
				client, err := elastic.NewClient(elastic.SetErrorLog(errorlog),elastic.SetSniff(false), elastic.SetURL(opt.DSN),elastic.SetBasicAuth(opt.USER,opt.PASSWORD),elastic.SetHttpClient( &http.Client{
					Transport:&http.Transport{
						DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
							var dial = &net.Dialer{
								Timeout:10*time.Second,
								KeepAlive:20*time.Minute,
							}
							return dial.DialContext(ctx,network,addr)
						},
						TLSClientConfig: &tls.Config{
							MaxVersion:         tls.VersionTLS11,
							InsecureSkipVerify: true,
						},
						MaxIdleConnsPerHost:100,
						MaxConnsPerHost:100,
						Dial: func(network, addr string) (conn net.Conn, e error) {
							var dial = &net.Dialer{
								Timeout:10*time.Second,
								KeepAlive:20*time.Minute,
							}
							return dial.Dial(network,addr)
						},
					},
				}))
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
