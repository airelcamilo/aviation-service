package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"aviation-service/internal/dto"
	"aviation-service/internal/service"
)

type WeatherHandler struct {
	logger  *zap.SugaredLogger
	service service.IWeatherService
}

func NewWeatherHandler(logger *zap.SugaredLogger, service service.IWeatherService) *WeatherHandler {
	return &WeatherHandler{
		logger:  logger,
		service: service,
	}
}

func (h *WeatherHandler) RegisterRoutes(r chi.Router) {
	r.Route("/weather", func(r chi.Router) {
		r.Get("/", h.GetWeather)
	})
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")

	weather, err := h.service.GetWeather(r.Context(), city)
	if err != nil {
		h.logger.Errorw("Failed to get weather: %s", err)
		respondWithError(w, http.StatusBadRequest, "Failed to get weather")
		return
	}

	h.logger.Info("Weather data get successfully")
	if weather == nil {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(nil, "No weather found"))
	} else {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(weather, ""))
	}
}
