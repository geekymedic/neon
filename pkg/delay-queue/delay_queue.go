package delay_queue

import (
	"strings"
	"sync"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/shamaton/msgpack"

	"github.com/geekymedic/neon/logger"
)

type TaskHandle = func(task Task)

type DelayQueue interface {
	Watch(handle TaskHandle, topic ...string)
	Put(topic string, priority int, delay, ttl time.Duration, task Task) error
	Delete(topic string, id uint64) error
	Update(task Task) error
}

type Task struct {
	Id   uint64
	Body []byte
}

type BeanstalkDelayQueue struct {
	once         sync.Once
	conn         *beanstalk.Conn
	watchConn    *beanstalk.Conn
	watchTimeout time.Duration
	logger       logger.Logger
}

func (delayQueue *BeanstalkDelayQueue) Watch(handle TaskHandle, topic ...string) {
	delayQueue.once.Do(func() {
		var _topics = strings.Join(topic, "-")
		log := delayQueue.logger.With("topic", _topics)
		tubeSet := beanstalk.NewTubeSet(delayQueue.watchConn, topic...)
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
	})
}

func (delayQueue *BeanstalkDelayQueue) Put(topic string, priority uint32, delay, ttr time.Duration, task Task) error {
	tube := beanstalk.Tube{
		Conn: delayQueue.conn,
		Name: topic,
	}
	buf, err := msgpack.Encode(task)
	if err != nil {
		return err
	}
	_, err = tube.Put(buf, priority, delay, ttr)
	return err
}

func (delayQueue *BeanstalkDelayQueue) Delete(_topic string, id uint64) error {
	return delayQueue.conn.Delete(id)
}

func (delayQueue *BeanstalkDelayQueue) Update(_topic string, id uint64, priority uint32, delay time.Duration) error {
	return delayQueue.conn.Release(id, priority, delay)
}
