package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MockRedis struct {
	Store map[string]string
}

func (m *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	val, ok := m.Store[key]
	result := redis.NewStringResult(val, nil)
	if !ok {
		result = redis.NewStringResult("", redis.Nil)
	}
	return result
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.Store[key] = value.(string)
	return redis.NewStatusResult("OK", nil)
}

func (m *MockRedis) Close() error {
	return nil
}

type MockRedisSetError struct{}

func (m *MockRedisSetError) Get(ctx context.Context, key string) *redis.StringCmd {
	return redis.NewStringResult("", redis.Nil)
}

func (m *MockRedisSetError) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusResult("", fmt.Errorf("Cache set failed"))
}

func (m *MockRedisSetError) Close() error { return nil }
