package services

import (
	"errors"
	"sync"
)

type Map struct {
	store map[string]string
	sync.RWMutex
}

// possible errors in map module
var (
	ErrKeyNotFound = errors.New("key not found")
)

func NewMap() *Map {
	return &Map{
		store: make(map[string]string),
	}
}

func (m *Map) Put(key, value string) {
	m.Lock()
	defer m.Unlock()
	m.store[key] = value
}

func (m *Map) Get(key string) (string, error) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.store[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	return value, nil
}
