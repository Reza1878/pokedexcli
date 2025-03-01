package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type cache struct {
	entry    map[string]cacheEntry
	mu       *sync.Mutex
	interval time.Duration
}

func (c *cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.entry[key]

	if !ok {
		return nil, ok
	}

	return v.Val, ok
}

func (c *cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		for k, v := range c.entry {
			elapsed := time.Since(v.CreatedAt)
			if elapsed > c.interval {
				delete(c.entry, k)
			}
		}
	}
}

func NewCache(interval time.Duration) *cache {
	c := cache{interval: interval}

	c.reapLoop()

	return &c
}
