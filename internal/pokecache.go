package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	entry    map[string]cacheEntry
	mu       *sync.Mutex
	interval time.Duration
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.entry[key]

	if !ok {
		return nil, ok
	}

	return v.Val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.entry {
			elapsed := time.Since(v.CreatedAt)
			if elapsed > c.interval {
				delete(c.entry, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{interval: interval, mu: &sync.Mutex{}, entry: map[string]cacheEntry{}}

	go c.reapLoop()

	return &c
}
