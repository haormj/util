package ttl

import (
	"time"
)

type Options struct {
	CleanupInterval time.Duration
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	option := Options{
		CleanupInterval: time.Minute * 10,
	}

	for _, o := range opts {
		o(&option)
	}

	return option
}

func CleanupInterval(cleanupInterval time.Duration) Option {
	return func(o *Options) {
		o.CleanupInterval = cleanupInterval
	}
}
