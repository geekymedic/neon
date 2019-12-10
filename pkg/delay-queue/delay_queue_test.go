package delay_queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/geekymedic/neon/utils/tool"
)

func TestNewBeanstalkDelayQueue(t *testing.T) {
	delayQueue, err := NewBeanstalkDelayQueue("127.0.0.1:11300")
	require.Nil(t, err)
	defer delayQueue.Close()
	delayQueue.Put("test", 1, time.Second*1, time.Second*2, Task{
		Id: tool.RandomUint64(),
		Body: []byte("{\"oss\": \"http://localhost.common\"}"),
	})
	delayQueue.Watch(func(task Task) {
		fmt.Println("taskId:", task.Id, string(task.Body))
	}, "test")
	time.Sleep(time.Second * 10)
}
