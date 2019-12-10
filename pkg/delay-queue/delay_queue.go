package delay_queue

import (
	"strings"
	"sync"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/shamaton/msgpack"

	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
)

type TaskHandle = func(task Task)

type DelayQueue interface {
	Watch(handle TaskHandle, topic ...string)
	Put(topic string, priority uint32, delay, ttr time.Duration, task Task) error
	Delete(topic string, id uint64) error
	Update(topic string, id uint64, priority uint32, delay time.Duration) error
	Close() error
}

type Task struct {
	Id   uint64
	Body []byte
}

type Conn struct {
	*beanstalk.Conn
}

type BeanstalkDelayQueue struct {
	once         *sync.Once
	conn         *Conn
	watchConn    *Conn
	watchTimeout time.Duration
	logger       logger.Logger
}

func NewBeanstalkDelayQueue(addr string) (*BeanstalkDelayQueue, error) {
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, errors.By(err)
	}
	watchConn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	delayQueue := &BeanstalkDelayQueue{
		once:      &sync.Once{},
		conn:      &Conn{conn},
		watchConn: &Conn{watchConn},
		logger:    logger.With("addr", addr),
	}
	return delayQueue, nil
}

func (delayQueue *BeanstalkDelayQueue) Watch(handle TaskHandle, topic ...string) {
	delayQueue.once.Do(func() {
		var _topics = strings.Join(topic, "-")
		log := delayQueue.logger.With("topic", _topics)
		tubeSet := beanstalk.NewTubeSet(delayQueue.watchConn.Conn, topic...)
		go func() {
			for {
				id, body, err := tubeSet.Reserve(delayQueue.watchTimeout)
				if err != nil {
					connErr, ok := err.(beanstalk.ConnError)
					if !ok {
						log.Error(err)
						continue
					}
					if connErr.Op != "reserve-with-timeout" {
						log.Error(connErr)
						continue
					}
					log.Debug("reserve timeout")
					continue
				}
				var task Task
				err = msgpack.Decode(body, &task)
				if err != nil {
					log.With("id", id).Error(err)
					continue
				}
				handle(task)
			}
		}()
	})
}

func (delayQueue *BeanstalkDelayQueue) Put(topic string, priority uint32, delay, ttr time.Duration, task Task) error {
	tube := beanstalk.Tube{
		Conn: delayQueue.conn.Conn,
		Name: topic,
	}
	buf, err := msgpack.Encode(task)
	if err != nil {
		return err
	}
	_, err = tube.Put(buf, priority, delay, ttr)
	return err
}

func (delayQueue *BeanstalkDelayQueue) Delete(topic string, id uint64) error {
	err := delayQueue.conn.Use(&beanstalk.Tube{Name: topic})
	if err != nil {
		return err
	}
	return delayQueue.conn.Delete(id)
}

func (delayQueue *BeanstalkDelayQueue) Update(topic string, id uint64, priority uint32, delay time.Duration) error {
	err := delayQueue.conn.Use(&beanstalk.Tube{Name: topic})
	if err != nil {
		return err
	}
	return delayQueue.conn.Release(id, priority, delay)
}

func(delayQueue *BeanstalkDelayQueue) Close() error {
	var err error
	delayQueue.once.Do(func(){
		err = delayQueue.conn.Close()
		if err != nil {
			return
		}
		err = delayQueue.watchConn.Close()
		return
	})

	return err
}