package cache

type MemoryOptions struct {
	MaxEntries int
}

type MemoryOption func(*MemoryOptions)

func newMemoryOptions(opts ...MemoryOption) MemoryOptions {
	option := MemoryOptions{
		MaxEntries: 1000,
	}

	for _, o := range opts {
		o(&option)
	}

	return option
}

func MaxEntries(maxEntries int) MemoryOption {
	return func(o *MemoryOptions) {
		o.MaxEntries = maxEntries
	}
}
