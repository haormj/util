package cache

// Cache define cache interface
type Cache interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Delete(string)
	Clear()
	GetString(string) (string, bool)
	GetStringE(string) (string, bool, error)
}

func NewCache() Cache {
	return NewMemoryCache()
}
