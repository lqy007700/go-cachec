package go_cachec

import (
	"context"
	"time"
)

// Cache 缓存接口定义
type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	//Del(key string) (any, error)
	LoadAndDel(ctx context.Context, key string) (any, error)
}
