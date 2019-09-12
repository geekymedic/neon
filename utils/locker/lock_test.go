package locker

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	t.Log(client)

	var key = fmt.Sprintf("%d", time.Now().Unix())
	var value = time.Now().Unix()
	var timeout = time.Second * 3
	t.Run("not held the lock", func(t *testing.T) {
		err := Unlock(client, key, value)
		if err != ErrNotHeld {
			t.Fatal(err)
		}
	})

	t.Run("get the lock", func(t *testing.T) {
		ok, err := Lock(client, key, value, timeout)
		if err != nil {
			t.Fatal(err)
		}
		if ok == false {
			t.Fatal(false)
		}

		ok, err = Lock(client, key, value, timeout)
		if err != nil {
			t.Fatal(err)
		}
		if ok == true {
			t.Fatal(err)
		}

		time.Sleep(timeout)
		ok, err = Lock(client, key, value, timeout)
		if err != nil {
			t.Fatal(err)
		}
		if ok == false {
			t.Fatal(false)
		}
	})
}
