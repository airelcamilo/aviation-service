package service_test

import (
	"aviation-service/config"
	"aviation-service/internal/dto"
	. "aviation-service/internal/mock"
	"aviation-service/internal/repository"
	. "aviation-service/internal/service"
	"aviation-service/pkg/logger"
	r "aviation-service/pkg/redis"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestAirportService_GetAllAirport(t *testing.T) {
	tests := []struct {
		name           string
		repo           repository.IAirportRepository
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Success get all airports",
			repo: &IAirportRepositoryMock{
				GetAllFunc: func(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX"}, {ID: 2, ICAO: "KAVL"}}, nil
				},
			},
			expectedResult: []dto.Airport{{ID: 1, ICAO: "KLAX"}, {ID: 2, ICAO: "KAVL"}},
		},
		{
			name: "Error get all airports",
			repo: &IAirportRepositoryMock{
				GetAllFunc: func(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
					return nil, fmt.Errorf("Failed to get all airports")
				},
			},
			expectedResult: ([]dto.Airport)(nil),
			expectedErr:    fmt.Errorf("Failed to get all airports"),
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{AIRPORT_API_URL: "http://123"}
			client := http.DefaultClient
			redisClient := &MockRedis{Store: make(map[string]string)}
			s := NewAirportService(log, tt.repo, cfg, client, redisClient)
			got, err := s.GetAllAirport(context.Background(), 20, 0)

			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, got)
			}
		})
	}
}

func TestAirportService_GetAirport(t *testing.T) {
	tests := []struct {
		name           string
		repo           repository.IAirportRepository
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Success get airport",
			repo: &IAirportRepositoryMock{
				GetByIdFunc: func(ctx context.Context, id int) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KLAX"}, nil
				},
			},
			expectedResult: &dto.Airport{ID: 1, ICAO: "KLAX"},
		},
		{
			name: "Error get airport",
			repo: &IAirportRepositoryMock{
				GetByIdFunc: func(ctx context.Context, id int) (*dto.Airport, error) {
					return nil, fmt.Errorf("Failed to get airport")
				},
			},
			expectedResult: (*dto.Airport)(nil),
			expectedErr:    fmt.Errorf("Failed to get airport"),
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{AIRPORT_API_URL: "http://123"}
			client := http.DefaultClient
			redisClient := &MockRedis{Store: make(map[string]string)}
			s := NewAirportService(log, tt.repo, cfg, client, redisClient)
			got, err := s.GetAirport(context.Background(), 1)

			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, got)
			}
		})
	}
}

func TestAirportService_SearchAirport(t *testing.T) {
	tests := []struct {
		name           string
		repo           *IAirportRepositoryMock
		httpClient     *mockHTTPClient
		redisClient    r.RedisClient
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Success search airport from cache",
			redisClient: &MockRedis{Store: map[string]string{
				"airport:KLAX:Lorem Ipsum": `[{"id": 1, "icao_ident": "KLAX"}]`,
			}},
			expectedResult: []dto.Airport{{ID: 1, ICAO: "KLAX"}},
		},
		{
			name: "Success search airport from repo",
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX"}}, nil
				},
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: []dto.Airport{{ID: 1, ICAO: "KLAX"}},
		},
		{
			name: "Success search airport from API",
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{}, nil
				},
				InsertFunc: func(ctx context.Context, airport *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KLAX"}, nil
				},
			},
			httpClient: &mockHTTPClient{
				response: `{"KLAX": [{"id": 1, "icao_ident": "KLAX"}]}`,
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: []dto.Airport{{ID: 1, ICAO: "KLAX"}},
		},
		{
			name: "Success search no airport found from API",
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{}, nil
				},
				InsertFunc: func(ctx context.Context, airport *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KLAX"}, nil
				},
			},
			httpClient: &mockHTTPClient{
				response: `{"KLAX": []}`,
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: ([]dto.Airport)(nil),
		},
		{
			name: "Error fetching airports data",
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{}, nil
				},
			},
			httpClient: &mockHTTPClient{
				err: fmt.Errorf("Error fetching airports data"),
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: ([]dto.Airport)(nil),
			expectedErr:    fmt.Errorf("Error fetching airports data"),
		},
		{
			name: "Error decoding body",
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{}, nil
				},
			},
			httpClient: &mockHTTPClient{
				response: `A`,
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: ([]dto.Airport)(nil),
			expectedErr:    fmt.Errorf("invalid character 'A' looking for beginning of value"),
		},
		{
			name: "Error insert airport from API",
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{}, nil
				},
				InsertFunc: func(ctx context.Context, airport *dto.Airport) (*dto.Airport, error) {
					return nil, fmt.Errorf("Error insert airport")
				},
			},
			httpClient: &mockHTTPClient{
				response: `{"KLAX": [{"id": 1, "icao_ident": "KLAX"}]}`,
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: ([]dto.Airport)(nil),
			expectedErr:    fmt.Errorf("Error insert airport"),
		},
		{
			name: "Error search airport",
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return nil, fmt.Errorf("Failed to search airport")
				},
			},
			redisClient: &MockRedis{Store: make(map[string]string)},
			expectedResult: ([]dto.Airport)(nil),
			expectedErr:    fmt.Errorf("Failed to search airport"),
		},
		{
			name: "Error set to cache (data from repo)",
			redisClient: &MockRedisSetError{},
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX"}}, nil
				},
			},
			expectedResult: []dto.Airport{{ID: 1, ICAO: "KLAX"}},
		},
		{
			name: "Error set to cache (data from API)",
			redisClient: &MockRedisSetError{},
			repo: &IAirportRepositoryMock{
				GetByICAOOrFacilityNameFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{}, nil
				},
				InsertFunc: func(ctx context.Context, airport *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KLAX"}, nil
				},
			},
			httpClient: &mockHTTPClient{
				response: `{"KLAX": [{"id": 1, "icao_ident": "KLAX"}]}`,
			},
			expectedResult: []dto.Airport{{ID: 1, ICAO: "KLAX"}},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{AIRPORT_API_URL: "http://123"}
			s := NewAirportService(log, tt.repo, cfg, tt.httpClient, tt.redisClient)
			got, err := s.SearchAirport(context.Background(), "KLAX", "Lorem Ipsum", 20, 0)

			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, got)
			}
		})
	}
}

func TestAirportService_CreateAirport(t *testing.T) {
	tests := []struct {
		name           string
		repo           *IAirportRepositoryMock
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Success create airport",
			repo: &IAirportRepositoryMock{
				InsertFunc: func(ctx context.Context, request *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KLAX"}, nil
				},
			},
			expectedResult: &dto.Airport{ID: 1, ICAO: "KLAX"},
		},
		{
			name: "Error create airport",
			repo: &IAirportRepositoryMock{
				InsertFunc: func(ctx context.Context, request *dto.Airport) (*dto.Airport, error) {
					return nil, fmt.Errorf("Failed to create airport")
				},
			},
			expectedResult: (*dto.Airport)(nil),
			expectedErr:    fmt.Errorf("Failed to create airport"),
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{AIRPORT_API_URL: "http://123"}
			client := http.DefaultClient
			redisClient := &MockRedis{Store: make(map[string]string)}
			s := NewAirportService(log, tt.repo, cfg, client, redisClient)
			got, err := s.CreateAirport(context.Background(), &dto.Airport{})

			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, got)
			}
		})
	}
}

func TestAirportService_UpdateAirport(t *testing.T) {
	tests := []struct {
		name           string
		repo           *IAirportRepositoryMock
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Success update airport",
			repo: &IAirportRepositoryMock{
				UpdateByIdFunc: func(ctx context.Context, request *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KLAX"}, nil
				},
			},
			expectedResult: &dto.Airport{ID: 1, ICAO: "KLAX"},
		},
		{
			name: "Error update airport",
			repo: &IAirportRepositoryMock{
				UpdateByIdFunc: func(ctx context.Context, request *dto.Airport) (*dto.Airport, error) {
					return nil, fmt.Errorf("Failed to update airport")
				},
			},
			expectedResult: (*dto.Airport)(nil),
			expectedErr:    fmt.Errorf("Failed to update airport"),
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{AIRPORT_API_URL: "http://123"}
			client := http.DefaultClient
			redisClient := &MockRedis{Store: make(map[string]string)}
			s := NewAirportService(log, tt.repo, cfg, client, redisClient)
			got, err := s.UpdateAirport(context.Background(), &dto.Airport{})

			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, got)
			}
		})
	}
}

func TestAirportService_DeleteAirport(t *testing.T) {
	tests := []struct {
		name           string
		repo           *IAirportRepositoryMock
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Success delete airport",
			repo: &IAirportRepositoryMock{
				DeleteFunc: func(ctx context.Context, id int) error {
					return nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error delete airport",
			repo: &IAirportRepositoryMock{
				DeleteFunc: func(ctx context.Context, id int) error {
					return fmt.Errorf("Failed to create airport")
				},
			},
			expectedErr: fmt.Errorf("Failed to create airport"),
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{AIRPORT_API_URL: "http://123"}
			client := http.DefaultClient
			redisClient := &MockRedis{Store: make(map[string]string)}
			s := NewAirportService(log, tt.repo, cfg, client, redisClient)
			err := s.DeleteAirport(context.Background(), 1)

			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

type mockHTTPClient struct {
	response string
	body     io.ReadCloser
	err      error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.body != nil {
		return &http.Response{
			StatusCode: 200,
			Body:       m.body,
		}, nil
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(m.response)),
	}, nil
}
