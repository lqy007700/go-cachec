package go_cachec

import (
	"context"
	"time"
)

type ReadThroughCache struct {
	Cache
	LoadFunc   func(ctx context.Context, key string) (any, error)
	expiration time.Duration
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	// 无数据
	if err == ErrKeyNotFund {
		// 方法存在
		if r.LoadFunc != nil {
			// 查询数据
			val, err = r.LoadFunc(ctx, key)
			// 写入cache
			_ = r.Cache.Set(ctx, key, val, r.expiration)
		}
	}

	return val, nil
}
