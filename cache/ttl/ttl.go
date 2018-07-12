package ttl

import (
	"time"

	pc "github.com/patrickmn/go-cache"
)

type TtlCache struct {
	cache  *pc.Cache
	option Options
}

func NewTtlCache(opts ...Option) *TtlCache {
	option := newOptions(opts...)

	ttlCache := &TtlCache{
		option: option,
		cache:  pc.New(0, option.CleanupInterval),
	}
	return ttlCache
}

func (tc *TtlCache) Set(key string, value interface{}, ttl time.Duration) {
	tc.cache.Set(key, value, ttl)

}

func (tc *TtlCache) Get(key string) (interface{}, bool) {
	return tc.cache.Get(key)
}

func (tc *TtlCache) Delete(key string) {
	tc.cache.Delete(key)
}

func (tc *TtlCache) Clear() {
	tc.cache.Flush()
}
