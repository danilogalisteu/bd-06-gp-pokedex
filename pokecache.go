package main

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
	c.cache[key] = cacheEntry{time.Now(), val}
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
