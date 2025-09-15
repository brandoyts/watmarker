package cache

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

type RedisClientConfig struct {
	Addr     string
	Username string
	Password string
}

func NewRedisClient(config RedisClientConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:      config.Addr,
		Username:  config.Username,
		Password:  config.Password,
		TLSConfig: &tls.Config{},
	})

	return &RedisCache{
		client: client,
	}, nil
}

func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

func (rc *RedisCache) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}
