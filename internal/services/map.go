package services

import (
	"gophermap/internal/custom_error"
	"sync"
)

type Map struct {
	store map[string]string
	sync.RWMutex
}

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
		return "", &custom_error.AppError{
			Code:    404,
			Message: "key not found",
		}
	}

	return value, nil
}

func (m *Map) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	_, ok := m.store[key]
	if !ok {
		return &custom_error.AppError{
			Code:    404,
			Message: "key not found",
		}
	}

	delete(m.store, key)
	return nil
}
