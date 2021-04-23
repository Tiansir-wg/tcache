package main

import (
	"sync"
	"tcache/lru"
)

type cache struct {
	mu           sync.Mutex
	lru          *lru.Cache
	limitedBytes int64
}

func (c *cache) put(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 采用延迟初始化方式
	if c.lru == nil {
		c.lru = lru.New(c.limitedBytes, nil)
	}
	c.lru.Put(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 未初始化
	if c.lru == nil {
		return
	}
	// 该值在缓存中
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
