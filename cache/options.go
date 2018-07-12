package cache

import (
	"github.com/haormj/util/cache/ttl"
)

type Options struct {
	BaseCache BaseCache
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	option := Options{
		BaseCache: ttl.NewTtlCache(),
	}

	for _, o := range opts {
		o(&option)
	}

	return option
}
