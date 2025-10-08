package utils_test

import (
	. "aviation-service/internal/mock"
	"aviation-service/internal/utils"
	r "aviation-service/pkg/redis"
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetStruct(t *testing.T) {
	type Airport struct {
		ID   int    `json:"id"`
		ICAO string `json:"icao_ident"`
	}

	ctx := context.Background()

	tests := []struct {
		name           string
		redisClient    r.RedisClient
		key            string
		expectedErr    error
		expectedResult []Airport
	}{
		{
			name: "Success Get Struct from cache",
			redisClient: &MockRedis{Store: map[string]string{
				"airport:KLAX:": `[{"id":1,"icao_ident":"KLAX"}]`,
			}},
			key:            "airport:KLAX:",
			expectedResult: []Airport{{ID: 1, ICAO: "KLAX"}},
		},
		{
			name:        "Error Get Struct failed",
			redisClient: &MockRedis{Store: map[string]string{}},
			key:         "airport:KLAX:",
			expectedErr: fmt.Errorf("redis: nil"),
		},
		{
			name: "Error Get Struct unmarshal",
			redisClient: &MockRedis{Store: map[string]string{
				"airport:KLAX:": "not-json",
			}},
			key:         "airport:KLAX:",
			expectedErr: fmt.Errorf("invalid character 'o' in literal null (expecting 'u')"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []Airport
			err := utils.GetStruct(tt.redisClient, ctx, tt.key, &result)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, result)
			}
		})
	}
}

func TestSetStruct(t *testing.T) {
	type Airport struct {
		ID   int    `json:"id"`
		ICAO string `json:"icao_ident"`
	}

	ctx := context.Background()

	tests := []struct {
		name         string
		redisClient  r.RedisClient
		key          string
		input        interface{}
		expectedErr  error
		expectedData string
	}{
		{
			name:         "Success Set Struct to cache",
			redisClient:  &MockRedis{Store: make(map[string]string)},
			key:          "airport:KLAX:",
			input:        []Airport{{ID: 1, ICAO: "KLAX"}},
			expectedData: `[{"id":1,"icao_ident":"KLAX"}]`,
		},
		{
			name:        "Error Set Struct json marshal",
			redisClient: &MockRedis{Store: make(map[string]string)},
			key:         "airport:KLAX:",
			input:       make(chan int),
			expectedErr: fmt.Errorf("json: unsupported type: chan int"),
		},
		{
			name:        "Error Set Struct failed",
			redisClient: &MockRedisSetError{},
			key:         "airport:KLAX:",
			input:       []Airport{{ID: 1, ICAO: "KLAX"}},
			expectedErr: fmt.Errorf("Cache set failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.SetStruct(tt.redisClient, ctx, tt.key, tt.input, time.Hour)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
			got, _ := tt.redisClient.Get(context.Background(), tt.key).Result()
			if got != tt.expectedData {
				t.Errorf("Expected data %v, got %v", tt.expectedData, got)
			}
		})
	}
}
