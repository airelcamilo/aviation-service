package service_test

import (
	"aviation-service/config"
	"aviation-service/internal/dto"
	. "aviation-service/internal/service"
	. "aviation-service/internal/mock"
	"aviation-service/pkg/logger"
	r "aviation-service/pkg/redis"
	"context"
	"fmt"
	"reflect"
	"testing"
)

type failingReadCloser struct{}

func (f *failingReadCloser) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("Read error")
}

func (f *failingReadCloser) Close() error {
	return nil
}


func TestWeatherService_GetWeather(t *testing.T) {
	tests := []struct {
		name           string
		httpClient     *mockHTTPClient
		redisClient    r.RedisClient
		expectedResult *dto.Weather
		expectedErr    error
	}{
		{
			name: "Success with data (cache miss)",
			httpClient: &mockHTTPClient{
				response: `{"location":{"name":"Asheville"},"current":{"last_updated":"2025-09-29 02:45","temp_c":17.2,"is_day":0}}`,
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: &dto.Weather{
				LastUpdated: "2025-09-29 02:45",
				TempC:       17.2,
				IsDay:       0,
			},
		},
		{
			name:       "Success with data (cache hit)",
			httpClient: &mockHTTPClient{},
			redisClient: &MockRedis{Store: map[string]string{
				"weather:Asheville": `{"current":{"last_updated":"2025-09-29 02:45","temp_c":17.2,"is_day":0}}`,
			}},
			expectedResult: &dto.Weather{
				LastUpdated: "2025-09-29 02:45",
				TempC:       17.2,
				IsDay:       0,
			},
		},
		{
			name: "Error get from cache",
			httpClient: &mockHTTPClient{
				response: `{"location":{"name":"Asheville"},"current":{"last_updated":"2025-09-29 02:45","temp_c":17.2,"is_day":0}}`,
			},
			redisClient: &MockRedis{Store: map[string]string{
				"weather:Washington": `A`,
			}},
			expectedResult: &dto.Weather{
				LastUpdated: "2025-09-29 02:45",
				TempC:       17.2,
				IsDay:       0,
			},
		},
		{
			name: "Error reading response body",
			httpClient: &mockHTTPClient{
				body: &failingReadCloser{},
			},
			redisClient: &MockRedis{Store: map[string]string{}},
			expectedErr: fmt.Errorf("Read error"),
		},
		{
			name: "Error fetching weather data",
			redisClient: &MockRedis{Store: make(map[string]string)},
			httpClient: &mockHTTPClient{
				err: fmt.Errorf("Error fetching weather data"),
			},
			expectedErr: fmt.Errorf("Error fetching weather data"),
		},
		{
			name: "Error decoding body",
			redisClient: &MockRedis{Store: make(map[string]string)},
			httpClient: &mockHTTPClient{
				response: `A`,
			},
			expectedErr: fmt.Errorf("invalid character 'A' looking for beginning of value"),
		},
		{
			name: "Error set to cache",
			httpClient: &mockHTTPClient{
				response: `{"location":{"name":"Asheville"},"current":{"last_updated":"2025-09-29 02:45","temp_c":17.2,"is_day":0}}`,
			},
			redisClient: &MockRedisSetError{},
			expectedResult: &dto.Weather{
				LastUpdated: "2025-09-29 02:45",
				TempC:       17.2,
				IsDay:       0,
			},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{AIRPORT_API_URL: "http://123"}
			s := NewWeatherService(log, cfg, tt.httpClient, tt.redisClient)

			got, err := s.GetWeather(context.Background(), "Asheville")
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, got)
			}
		})
	}
}
