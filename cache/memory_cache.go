package cache

import (
	"github.com/golang/groupcache/lru"
	"github.com/spf13/cast"
)

// MemoryCache implements Cache by using memory
type MemoryCache struct {
	option MemoryOptions
	cache  *lru.Cache
}

func NewMemoryCache(opts ...MemoryOption) *MemoryCache {
	option := newMemoryOptions(opts...)

	cache := lru.New(option.MaxEntries)

	memoryCache := &MemoryCache{
		option: option,
		cache:  cache,
	}

	return memoryCache
}

func (mc *MemoryCache) Set(key string, value interface{}) {
	mc.cache.Add(key, value)
}

func (mc *MemoryCache) Delete(key string) {
	mc.cache.Remove(key)
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
	return mc.cache.Get(key)
}

func (mc *MemoryCache) GetString(key string) (string, bool) {
	value, ok := mc.Get(key)
	v := cast.ToString(value)
	return v, ok
}
func (mc *MemoryCache) GetStringE(key string) (string, bool, error) {
	value, ok := mc.Get(key)
	vv, err := cast.ToStringE(value)
	return vv, ok, err
}

func (mc *MemoryCache) Clear() {
	mc.cache.Clear()
}
