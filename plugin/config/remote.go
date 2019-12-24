package config

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"

	"github.com/geekymedic/neon/errors"
)

type Event struct {
	EV interface{}
}

type Backend interface {
	Get(ctx context.Context, path string) ([]byte, error)
	Watch(ctx context.Context, path string) <-chan Event
	Close() error
}

// Only support etcd v3.0
type remoteConfigProvider struct {
	onlyOnce *sync.Once
	etcd     Backend
}

func (rc *remoteConfigProvider) initBackend(rp viper.RemoteProvider) {
	if rp.Provider() == "etcd" {
		rc.onlyOnce.Do(func() {
			var (
				endpoint           = rp.Endpoint()
				secretKeyString    = strings.SplitN(rp.SecretKeyring(), "@", 2)
				userName, password string
				err                error
			)
			if len(secretKeyString) == 2 {
				userName = secretKeyString[0]
				password = secretKeyString[1]
			}
			rc.etcd, err = newEtcdBackend(userName, password, []string{endpoint}, time.Second*3)
			if err != nil {
				panic(err)
			}
			fmt.Println("init etcd provider")
		})
	}

	if rp.Provider() == "abolo" {

	}
}

func (rc *remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	rc.initBackend(rp)
	buf, err := rc.etcd.Get(context.Background(), rp.Path())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return bytes.NewReader(buf), nil
}

func (rc *remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	rc.initBackend(rp)
	buf, err := rc.etcd.Get(context.Background(), rp.Path())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return bytes.NewReader(buf), nil
}

func (rc *remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	rc.initBackend(rp)
	quitwc := make(chan bool)
	viperResponseCh := make(chan *viper.RemoteResponse)

	var backend Backend

	if rp.Provider() == "etcd" {
		backend = rc.etcd
	}

	go func(vr <-chan *viper.RemoteResponse, quitwc <-chan bool) {
		defer backend.Close()
		for {
			select {
			case <-quitwc:
				return
			case <-backend.Watch(context.Background(), rp.Path()):
				fmt.Println("receive a config change event")
				buf, err := backend.Get(context.Background(), rp.Path())
				viperResponseCh <- &viper.RemoteResponse{
					Error: err,
					Value: buf,
				}
			}
		}
	}(viperResponseCh, quitwc)
	return viperResponseCh, quitwc
}

type etcdBackend struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cli        *clientv3.Client
	once       *sync.Once
	watchEvent chan Event
}

func newEtcdBackend(userName, password string, endpoints []string, dialTimeout time.Duration) (*etcdBackend, error) {
	ctx, cancelFn := context.WithCancel(context.Background())
	cli, err := clientv3.New(clientv3.Config{
		Username:    userName,
		Password:    password,
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		return nil, err
	}
	backend := &etcdBackend{
		ctx:        ctx,
		cancel:     cancelFn,
		cli:        cli,
		watchEvent: make(chan Event),
		once:       &sync.Once{},
	}
	return backend, nil
}

func (backend *etcdBackend) Get(ctx context.Context, path string) ([]byte, error) {
	resp, err := backend.cli.Get(ctx, path, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var buf = bytes.NewBuffer(nil)
	for _, ev := range resp.Kvs {
		buf.Write(ev.Value)
		buf.WriteString("\n")
	}
	return buf.Bytes(), nil
}

func (backend *etcdBackend) Watch(ctx context.Context, path string) <-chan Event {
	backend.once.Do(func() {
		go func() {
			watchChan := backend.cli.Watch(ctx, path, clientv3.WithPrefix())
			for {
				select {
				case <-backend.ctx.Done():
					return
				case ev := <-watchChan:
					fmt.Println(ev.CompactRevision)
					backend.watchEvent <- Event{
						EV: ev,
					}
				}
			}
		}()
	})
	return backend.watchEvent
}

func (backend *etcdBackend) Close() error {
	backend.cancel()
	return backend.cli.Close()
}

func init() {
	viper.RemoteConfig = &remoteConfigProvider{
		onlyOnce: &sync.Once{},
	}
}
