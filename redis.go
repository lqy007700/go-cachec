package go_cachec

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client redis.Cmdable
}

func (r *RedisCache) LoadAndDel(ctx context.Context, key string) (any, error) {
	return r.client.GetDel(ctx, key).Result()
}

func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	result, err := r.client.Set(ctx, key, val, expiration).Result()
	if err != nil {
		return err
	}

	if result != "OK" {
		return errors.New("not OK")
	}
	return nil
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
