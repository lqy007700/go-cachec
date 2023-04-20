package go_cachec

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestClient_TryLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	c := &Client{
		client: rdb,
	}

	lock, err := c.TryLock(context.Background(), "name", time.Second*50)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println(lock)
	time.Sleep(time.Second * 10)

	err = lock.UnLock(context.Background())
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("释放成功")
}
