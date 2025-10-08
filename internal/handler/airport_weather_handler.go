package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"aviation-service/internal/dto"
)

type AirportWeatherHandler struct {
	logger  *zap.SugaredLogger
	service AirportWeatherService
}

type AirportWeatherService interface {
	SearchAirportWeather(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.AirportWeather, error)
}

func NewAirportWeatherHandler(logger *zap.SugaredLogger, service AirportWeatherService) *AirportWeatherHandler {
	return &AirportWeatherHandler{
		logger:  logger,
		service: service,
	}
}

func (h *AirportWeatherHandler) RegisterRoutes(r chi.Router) {
	r.Route("/airport-weather", func(r chi.Router) {
		r.Get("/", h.SearchAirportWeather)
	})
}

func (h *AirportWeatherHandler) SearchAirportWeather(w http.ResponseWriter, r *http.Request) {
	icao := r.URL.Query().Get("icao")
	facilityName := r.URL.Query().Get("facilityName")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	airportWeathers, err := h.service.SearchAirportWeather(r.Context(), icao, facilityName, pageSize, offset)
	if err != nil {
		h.logger.Errorw("Failed to get airport and weather", "error", err)
		respondWithError(w, http.StatusBadRequest, "Failed to get airport and weather")
		return
	}

	h.logger.Info("Airport and weather data get successfully")
	if airportWeathers == nil {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(nil, "No airport and weather found"))
	} else {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(dto.PaginatedResponse{
			Page:     page,
			PageSize: pageSize,
			Data:     airportWeathers,
		}, ""))
	}
}
