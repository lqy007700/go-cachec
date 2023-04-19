package go_cachec

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrOverCapacity = errors.New("超出容量限制")
)

// MaxCntCache 控制缓存键值对数量
type MaxCntCache struct {
	*BuildInMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(b *BuildInMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCache: b,
		maxCnt:          maxCnt,
	}

	// 重写onEvicted del key 更新当前键值对数量
	origin := b.onEvicted
	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}
	return res
}

func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			return ErrOverCapacity
		}
		c.cnt++
	}

	return c.set(key, val, expiration)
}
