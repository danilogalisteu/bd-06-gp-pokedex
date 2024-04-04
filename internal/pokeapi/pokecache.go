package pokeapi

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]cacheEntry
	lock  *sync.Mutex
}

func (c Cache) Add(key string, val []byte) {
	c.lock.Lock()
	c.cache[key] = cacheEntry{createdAt: time.Now(), val: val}
	c.lock.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	entry, ok := c.cache[key]
	c.lock.Unlock()
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for now := range ticker.C {
		c.lock.Lock()
		for key, item := range c.cache {
			if now.Sub(item.createdAt) > interval {
				delete(c.cache, key)
			}
		}
		c.lock.Unlock()
	}
}

func NewCache(interval time.Duration) Cache {
	c := Cache{cache: map[string]cacheEntry{}, lock: &sync.Mutex{}}
	go c.reapLoop(interval)
	return c
}
