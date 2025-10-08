package service_test

import (
	"aviation-service/internal/dto"
	. "aviation-service/internal/service"
	. "aviation-service/internal/mock"
	"aviation-service/pkg/logger"
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestAirportWeatherService_SearchAirportWeather(t *testing.T) {
	city := "ASHEVILLE"
	tests := []struct {
		name           string
		airportService IAirportService
		weatherService IWeatherService
		expectedResult []dto.AirportWeather
		expectedErr    error
	}{
		{
			name: "Success with airport and weather data",
			airportService: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, name string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX", City: &city}}, nil
				},
			},
			weatherService: &IWeatherServiceMock{
				GetWeatherFunc: func(ctx context.Context, city string) (*dto.Weather, error) {
					return &dto.Weather{
						LastUpdated: "2025-09-29 02:45",
						TempC:       17.2,
						IsDay:       0,
					}, nil
				},
			},
			expectedResult: []dto.AirportWeather{{
				Airport: dto.Airport{ID: 1, ICAO: "KLAX", City: &city},
				Weather: &dto.Weather{
					LastUpdated: "2025-09-29 02:45",
					TempC:       17.2,
					IsDay:       0,
				},
			}},
		},
		{
			name: "Success with airport data (doesn't have city value)",
			airportService: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, name string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX"}}, nil
				},
			},
			weatherService: &IWeatherServiceMock{
				GetWeatherFunc: func(ctx context.Context, city string) (*dto.Weather, error) {
					return &dto.Weather{}, nil
				},
			},
			expectedResult: []dto.AirportWeather{{
				Airport: dto.Airport{ID: 1, ICAO: "KLAX"},
			}},
		},
		{
			name: "Success with airport data (doesn't have weather data)",
			airportService: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, name string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KLAX", City: &city}}, nil
				},
			},
			weatherService: &IWeatherServiceMock{
				GetWeatherFunc: func(ctx context.Context, city string) (*dto.Weather, error) {
					return nil, fmt.Errorf("Weather not available")
				},
			},
			expectedResult: []dto.AirportWeather{{
				Airport: dto.Airport{ID: 1, ICAO: "KLAX", City: &city},
			}},
		},
		{
			name: "Failed to search airports",
			airportService: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, name string, limit, offset int) ([]dto.Airport, error) {
					return nil, fmt.Errorf("Failed to search airports")
				},
			},
			weatherService: &IWeatherServiceMock{
				GetWeatherFunc: func(ctx context.Context, city string) (*dto.Weather, error) {
					return &dto.Weather{}, nil
				},
			},
			expectedErr: fmt.Errorf("Failed to search airports"),
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAirportWeatherService(log, tt.airportService, tt.weatherService)

			ctx := context.Background()

			got, err := s.SearchAirportWeather(ctx, "KLAX", "Lorem Ipsum", 20, 0)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Expected result %+v, got %+v", tt.expectedResult, got)
			}
		})
	}
}
