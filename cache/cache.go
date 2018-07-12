package cache

import (
	"time"

	"github.com/spf13/cast"
)

// Cache define cache interface
type Cache interface {
	BaseCache
	GetString(string) (string, bool)
	GetStringE(string) (string, bool, error)
}

type BaseCache interface {
	Set(string, interface{}, time.Duration)
	Get(string) (interface{}, bool)
	Delete(string)
	Clear()
}

func NewCache(opts ...Option) Cache {
	option := newOptions(opts...)

	defaultCache := defaultCache{
		BaseCache: option.BaseCache,
		option:    option,
	}

	return defaultCache
}

type defaultCache struct {
	BaseCache
	option Options
}

func (dc defaultCache) GetString(key string) (string, bool) {
	v, ok := dc.Get(key)
	value := cast.ToString(v)
	return value, ok
}

func (dc defaultCache) GetStringE(key string) (string, bool, error) {
	v, ok := dc.Get(key)
	value, err := cast.ToStringE(v)
	return value, ok, err
}
