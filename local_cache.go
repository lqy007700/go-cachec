package go_cachec

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFund = errors.New("key不存在")
)

type BuildInMapCache struct {
	data  map[string]*Item
	mu    sync.RWMutex
	close chan struct{}
}

type Item struct {
	val      any
	deadline time.Time
}

func NewBuildInMapCache(size int32) *BuildInMapCache {
	res := &BuildInMapCache{
		data:  make(map[string]*Item, size),
		close: make(chan struct{}, 1),
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)

		for {
			select {
			case t := <-ticker.C:
				res.mu.Lock()
				i := 0
				for key, val := range res.data {
					if i >= 1000 {
						break
					}
					if !val.deadline.IsZero() && val.deadline.Before(t) {
						delete(res.data, key)
					}
					i++
				}
				res.mu.Unlock()
			case <-res.close:
				return
			}
		}
	}()

	return res
}

func (b *BuildInMapCache) Get(ctx context.Context, key string) (any, error) {
	b.mu.RLock()
	v, ok := b.data[key]
	b.mu.RUnlock()
	if !ok {
		return nil, ErrKeyNotFund
	}

	now := time.Now()
	if !v.deadline.IsZero() && v.deadline.Before(now) {
		b.mu.Lock()
		v, ok = b.data[key]
		if ok && !v.deadline.IsZero() && v.deadline.Before(now) {
			delete(b.data, key)
		}
		b.mu.Unlock()
		return nil, ErrKeyNotFund
	}

	return v.val, nil
}

func (b *BuildInMapCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[key] = &Item{
		val:      val,
		deadline: time.Now().Add(expiration),
	}

	// set的时候开启一个定时任务，指定过期时间后删除key
	//if expiration > 0 {
	//	time.AfterFunc(expiration, func() {
	//		v, ok := b.data[key]
	//		if ok && v.deadline.Before(time.Now()) {
	//			b.Del(ctx, key)
	//		}
	//	})
	//}

	return nil
}

func (b *BuildInMapCache) Del(ctx context.Context, key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.data, key)
	return nil
}

func (b *BuildInMapCache) Close() error {
	select {
	case b.close <- struct{}{}:
	default:
		return errors.New("已关闭")
	}
	return nil
}
