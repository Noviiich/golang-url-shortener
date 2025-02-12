package mock

import (
	"context"
	"errors"
	"time"
)

type MockRedisCache struct {
	Repository map[string]string
	TTL        map[string]time.Time
}

func NewMockRedisCache() *MockRedisCache {
	return &MockRedisCache{
		Repository: make(map[string]string),
		TTL:        make(map[string]time.Time),
	}
}

func (m *MockRedisCache) Set(ctx context.Context, key, val string) (string, error) {
	m.Repository[key] = val
	m.TTL[key] = time.Now().Add(time.Minute)
	return key, nil
}

func (m *MockRedisCache) Get(ctx context.Context, key string) (string, error) {
	val, ok := m.Repository[key]
	if !ok {
		return "", errors.New("key now found")
	}

	if time.Now().After(m.TTL[key]) {
		delete(m.Repository, key)
		delete(m.TTL, key)
		return "", errors.New("key expired")
	}

	return val, nil
}

func (m *MockRedisCache) Delete(ctx context.Context, key string) error {
	_, ok := m.Repository[key]
	if !ok {
		return errors.New("key not found")
	}
	delete(m.Repository, key)
	delete(m.TTL, key)
	return nil
}
