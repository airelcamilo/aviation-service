package service_test

import (
	"aviation-service/internal/dto"
	. "aviation-service/internal/mock"
	"aviation-service/internal/repository"
	. "aviation-service/internal/service"
	"aviation-service/pkg/logger"
	"context"
	"fmt"
	"testing"
)

type expectedCount struct {
	total   int
	success int
	failed  int
	err     int
}

func TestAviationSyncService_Sync(t *testing.T) {
	tests := []struct {
		name           string
		repo           repository.IAirportRepository
		airportService IAirportService
		expected       expectedCount
		expectedErr    error
	}{
		{
			name: "Success with data",
			repo: &IAirportRepositoryMock{
				GetAllPendingFunc: func(ctx context.Context) ([]dto.Airport, error) {
					var airports []dto.Airport
					n := 30
					for i := 1; i <= n; i++ {
						airports = append(airports, dto.Airport{
							ID:   i,
							ICAO: fmt.Sprintf("KA%02d", i),
						})
					}
					return airports, nil
				},
				UpdateByICAOFunc: func(ctx context.Context, airport []dto.Airport) error {
					return nil
				},
			},
			airportService: &IAirportServiceMock{
				FetchAirportDataFunc: func(icao string) (*dto.AirportDataResponse, error) {
					return &dto.AirportDataResponse{
						"KLAX": []dto.Airport{}, 
						"KAVL": []dto.Airport{{ID: 1, ICAO: "KAVL"}},
						}, nil
				},
			},
			expected: expectedCount{
				total:   30,
				success: 1,
				failed:  1,
			},
		},
		{
			name: "Success no data to sync",
			repo: &IAirportRepositoryMock{
				GetAllPendingFunc: func(ctx context.Context) ([]dto.Airport, error) {
					return []dto.Airport{}, nil
				},
			},
		},
		{
			name: "Failed get all pending airports",
			repo: &IAirportRepositoryMock{
				GetAllPendingFunc: func(ctx context.Context) ([]dto.Airport, error) {
					return nil, fmt.Errorf("Failed get all pending airports")
				},
			},
			expectedErr: fmt.Errorf("Failed get all pending airports"),
		},
		{
			name: "Error fetching airport data from API",
			repo: &IAirportRepositoryMock{
				GetAllPendingFunc: func(ctx context.Context) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX"}}, nil
				},
				UpdateByICAOFunc: func(ctx context.Context, airport []dto.Airport) error {
					return fmt.Errorf("Error fetching airport data")
				},
			},
			airportService: &IAirportServiceMock{
				FetchAirportDataFunc: func(icao string) (*dto.AirportDataResponse, error) {
					return nil, fmt.Errorf("Error fetching airport data")
				},
			},
			expected: expectedCount{
				total: 1,
				err:   1,
			},
		},
		{
			name: "Error updating airport data",
			repo: &IAirportRepositoryMock{
				GetAllPendingFunc: func(ctx context.Context) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX"}}, nil
				},
				UpdateByICAOFunc: func(ctx context.Context, airport []dto.Airport) error {
					return fmt.Errorf("Error updating airport data")
				},
			},
			airportService: &IAirportServiceMock{
				FetchAirportDataFunc: func(icao string) (*dto.AirportDataResponse, error) {
					return &dto.AirportDataResponse{
						"KLAX": []dto.Airport{{ID: 1, ICAO: "KLAX"}},
						}, nil
				},
			},
			expected: expectedCount{
				total: 1,
				err:   1,
			},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAviationSyncService(log, tt.repo, tt.airportService)

			resp, err := s.Sync(context.Background())
			if err != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if resp == nil {
				return
			}

			if resp.Total != tt.expected.total {
				t.Errorf("Expected total %d, got %d", tt.expected.total, resp.Total)
			}
			if resp.Success != tt.expected.success {
				t.Errorf("Expected success %d, got %d", tt.expected.success, resp.Success)
			}
			if resp.Failed != tt.expected.failed {
				t.Errorf("Expected failed %d, got %d", tt.expected.failed, resp.Failed)
			}
			if resp.Error != tt.expected.err {
				t.Errorf("Expected failed %d, got %d", tt.expected.err, resp.Error)
			}
		})
	}
}
