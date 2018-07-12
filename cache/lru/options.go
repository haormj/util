package lru

type Options struct {
	MaxEntries int
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	option := Options{
		MaxEntries: 1000,
	}

	for _, o := range opts {
		o(&option)
	}

	return option
}

func MaxEntries(maxEntries int) Option {
	return func(o *Options) {
		o.MaxEntries = maxEntries
	}
}
