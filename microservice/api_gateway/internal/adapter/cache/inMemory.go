package cache

import (
	"context"
	"sync"
	"time"
)

type InMemoryCache struct {
	store map[string]int64
	mu    sync.Mutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{store: make(map[string]int64)}
}

func (c *InMemoryCache) Increment(ctx context.Context, key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key]++
	return c.store[key], nil
}

func (c *InMemoryCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}
