package util

import (
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
