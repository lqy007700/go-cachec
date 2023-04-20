package go_cachec

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/unlock.lua
var luaUnlock string

type Client struct {
	client redis.Cmdable
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {

	val := uuid.New().String()
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("抢锁失败")
	}
	return &Lock{
		client: c.client,
		key:    key,
		val:    val,
	}, nil
}

type Lock struct {
	client redis.Cmdable
	key    string
	val    string
}

func (l *Lock) UnLock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, []string{l.val}).Int64()
	if err != nil {
		return err
	}

	if res != 1 {
		return errors.New("释放失败")
	}
	return nil
}
