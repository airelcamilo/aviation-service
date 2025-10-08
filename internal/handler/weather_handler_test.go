package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"aviation-service/internal/dto"
	"aviation-service/internal/service"
	utils "aviation-service/internal/testutils"
	. "aviation-service/internal/handler"
	"aviation-service/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type mockWeatherService struct {
	weatherResponse *dto.Weather
	err             error
}

func (m *mockWeatherService) GetWeather(ctx context.Context, city string) (*dto.Weather, error) {
	return m.weatherResponse, m.err
}

func TestWeatherHandler_GetWeather(t *testing.T) {
	tests := []struct {
		name        string
		service     service.IWeatherService
		queryParams string
		utils.ExpectedResult
	}{
		{
			name:        "Success with data",
			service:     &mockWeatherService{weatherResponse: &dto.Weather{TempC: 40, IsDay: 0, WindKph: 29}},
			queryParams: "?city=ASHEVILLE",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data:   &dto.Weather{TempC: 40, IsDay: 0, WindKph: 29},
			},
		},
		{
			name:        "No data",
			service:     &mockWeatherService{weatherResponse: nil},
			queryParams: "?city=ASHEVILLE",
			ExpectedResult: utils.ExpectedResult{
				Status:  http.StatusOK,
				Message: "No weather found",
			},
		},
		{
			name:        "Service error",
			service:     &mockWeatherService{err: fmt.Errorf("DB error")},
			queryParams: "?city=ASHEVILLE",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to get weather",
			},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewWeatherHandler(log, tt.service)
			h.RegisterRoutes(r)

			req := httptest.NewRequest(http.MethodGet, "/weather"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			h.GetWeather(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}
