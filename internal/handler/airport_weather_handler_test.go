package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"aviation-service/internal/dto"
	. "aviation-service/internal/handler"
	utils "aviation-service/internal/testutils"
	"aviation-service/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type mockAirportWeatherService struct {
	response []dto.AirportWeather
	err      error
}

func (m *mockAirportWeatherService) SearchAirportWeather(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.AirportWeather, error) {
	return m.response, m.err
}

func TestAirportWeatherHandler_SearchAirportWeather(t *testing.T) {
	tests := []struct {
		name        string
		service     AirportWeatherService
		queryParams string
		utils.ExpectedResult
	}{
		{
			name: "Success with data and correct pagination",
			service: &mockAirportWeatherService{response: []dto.AirportWeather{{
				Airport: dto.Airport{ICAO: "KAVL"},
				Weather: &dto.Weather{TempC: 38.1}}}},
			queryParams: "?icao=KAVL&page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data: dto.PaginatedResponse{
					Page:     1,
					PageSize: 10,
					Data: []dto.AirportWeather{{
						Airport: dto.Airport{ICAO: "KAVL"},
						Weather: &dto.Weather{TempC: 38.1}}},
				},
			},
		},
		{
			name: "Success with data and incorrect pagination",
			service: &mockAirportWeatherService{response: []dto.AirportWeather{{
				Airport: dto.Airport{ICAO: "KAVL"},
				Weather: &dto.Weather{TempC: 38.1}}}},
			queryParams: "?icao=KAVL&page=-1&pageSize=A",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data: dto.PaginatedResponse{
					Page:     1,
					PageSize: 10,
					Data: []dto.AirportWeather{{
						Airport: dto.Airport{ICAO: "KAVL"},
						Weather: &dto.Weather{TempC: 38.1}}},
				},
			},
		},
		{
			name:        "No data",
			service:     &mockAirportWeatherService{response: nil},
			queryParams: "?icao=KAVL&page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status:  http.StatusOK,
				Message: "No airport and weather found",
			},
		},
		{
			name:        "Service error",
			service:     &mockAirportWeatherService{err: fmt.Errorf("DB error")},
			queryParams: "?icao=KAVL&page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to get airport and weather",
			},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAirportWeatherHandler(log, tt.service)
			h.RegisterRoutes(r)

			req := httptest.NewRequest(http.MethodGet, "/airport-weather"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			h.SearchAirportWeather(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}
