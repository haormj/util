package lru

import (
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

// LruCache implements Cache by using memory
type LruCache struct {
	rw     sync.RWMutex
	option Options
	cache  *lru.Cache
}

func NewLruCache(opts ...Option) *LruCache {
	option := newOptions(opts...)

	cache := lru.New(option.MaxEntries)

	lruCache := &LruCache{
		option: option,
		cache:  cache,
	}

	return lruCache
}

func (mc *LruCache) Set(key string, value interface{}, ttl time.Duration) {
	mc.rw.Lock()
	mc.cache.Add(key, value)
	mc.rw.Unlock()
}

func (mc *LruCache) Delete(key string) {
	mc.rw.Lock()
	mc.cache.Remove(key)
	mc.rw.Unlock()
}

func (mc *LruCache) Get(key string) (interface{}, bool) {
	mc.rw.RLock()
	defer mc.rw.RUnlock()
	return mc.cache.Get(key)
}

func (mc *LruCache) Clear() {
	mc.rw.Lock()
	mc.cache.Clear()
	mc.rw.Unlock()
}
