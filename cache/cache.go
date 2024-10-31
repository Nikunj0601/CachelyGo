package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

type Cache struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]CacheEntry)}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheEntry{Value: value, ExpiresAt: time.Now().Add(ttl)}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if !ok || item.ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	return item.Value, true
}
