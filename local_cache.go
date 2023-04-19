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
	data      map[string]*Item
	mu        sync.RWMutex
	close     chan struct{}
	onEvicted func(key string, val any)
}

func (b *BuildInMapCache) LoadAndDel(ctx context.Context, key string) (any, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	v, ok := b.data[key]
	if !ok {
		return nil, ErrOverCapacity
	}

	delete(b.data, key)
	return v.val, nil
}

type Item struct {
	val      any
	deadline time.Time
}

func NewBuildInMapCache(size int32, onEvicted func(key string, val any)) *BuildInMapCache {
	res := &BuildInMapCache{
		data:      make(map[string]*Item, size),
		close:     make(chan struct{}),
		onEvicted: onEvicted,
	}

	// 轮训删除过期key
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
					if val.deadlineBefore(t) {
						res.del(key)
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

	// 判断key是否过期
	now := time.Now()
	if v.deadlineBefore(now) {
		b.mu.Lock()
		v, ok = b.data[key]
		if ok && v.deadlineBefore(now) {
			b.del(key)
		}
		b.mu.Unlock()
		return nil, ErrKeyNotFund
	}

	return v.val, nil
}

func (b *BuildInMapCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// set的时候开启一个定时任务，指定过期时间后删除key
	//if expiration > 0 {
	//	time.AfterFunc(expiration, func() {
	//		v, ok := b.data[key]
	//		if ok && v.deadline.Before(time.Now()) {
	//			b.Del(ctx, key)
	//		}
	//	})
	//}

	return b.set(key, val, expiration)
}
func (b *BuildInMapCache) set(key string, val any, expiration time.Duration) error {
	b.data[key] = &Item{
		val:      val,
		deadline: time.Now().Add(expiration),
	}
	return nil
}

func (b *BuildInMapCache) Del(ctx context.Context, key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.del(key)
	return nil
}

func (b *BuildInMapCache) del(key string) {
	item, ok := b.data[key]
	if !ok {
		return
	}

	delete(b.data, key)
	b.onEvicted(key, item.val)
}

func (b *BuildInMapCache) Close() error {
	select {
	case b.close <- struct{}{}:
	default:
		return errors.New("已关闭")
	}
	return nil
}

func (i *Item) deadlineBefore(t time.Time) bool {
	return !i.deadline.IsZero() && i.deadline.Before(t)
}
