package port

import (
	"context"
	"time"
)

//go:generate mockgen -package mock -source=../port/cache.go -destination=../mock/cache_mock.go
type RateLimit interface {
	Ping(ctx context.Context) error
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, window time.Duration) error
}
