package cache

import (
	"log"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	cache.Set("hello", "world")
	v, ok := cache.GetString("hello")
	log.Println(v, ok)
}

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache(MaxItemSize(2))
	cache.Set("hello", "1")
	cache.Set("world", "2")
	cache.Set("nihao", "3")
	cache.Get("hello")
	log.Println(cache.Get("hello"))
}
