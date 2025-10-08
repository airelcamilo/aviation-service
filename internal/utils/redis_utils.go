package utils

import (
	"aviation-service/pkg/redis"
	"context"
	"encoding/json"
	"time"
)

func GetStruct(r redis.RedisClient, ctx context.Context, key string, out interface{}) error {
	cached, err := r.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	
	if jsonErr := json.Unmarshal([]byte(cached), out); jsonErr != nil {
		return jsonErr
	}
	return nil
}

func SetStruct(r redis.RedisClient, ctx context.Context, key string, in interface{}, expiration time.Duration) error {
	data, jsonErr := json.Marshal(in)
	if jsonErr != nil {
		return jsonErr
	}
	if err := r.Set(ctx, key, string(data), expiration).Err(); err != nil {
		return err
	}
	return nil
}