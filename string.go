package util

import (
	"strings"
	"sync"
)

type SecurityString struct {
	sync.RWMutex
	data string
}

func NewSecurityString() *SecurityString {
	return &SecurityString{}
}

func (ss *SecurityString) Set(data string) {
	ss.Lock()
	ss.data = data
	ss.Unlock()
}

func (ss *SecurityString) Get() (data string) {
	ss.RLock()
	data = ss.data
	ss.RUnlock()
	return
}

// TrimSpaceAndSlash trim space first, then trim slash
func TrimSpaceAndSlash(s string) string {
	s = strings.TrimSpace(s)
	return strings.Trim(s, "/")
}
