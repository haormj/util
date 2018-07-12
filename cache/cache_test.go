package cache

import (
	"log"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	cache.Set("hello", "world", time.Second*5)
	for {
		v, ok := cache.GetString("hello")
		if !ok {
			break
		}
		log.Println(v)
		time.Sleep(time.Second * 1)
	}
}
