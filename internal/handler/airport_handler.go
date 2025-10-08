package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"aviation-service/internal/dto"
	"aviation-service/internal/service"
)

type AirportHandler struct {
	logger    *zap.SugaredLogger
	service   service.IAirportService
	validator AirportValidator
}

type AirportValidator interface {
	IsComplete(req *dto.Airport) bool
	Validate(req *dto.Airport) error
}

func NewAirportHandler(logger *zap.SugaredLogger, service service.IAirportService, validator AirportValidator) *AirportHandler {
	return &AirportHandler{
		logger:    logger,
		service:   service,
		validator: validator,
	}
}

func (h *AirportHandler) RegisterRoutes(r chi.Router) {
	r.Route("/airport", func(r chi.Router) {
		r.Get("/", h.GetAllAirport)
		r.Post("/", h.CreateAirport)

		r.Route("/search", func(r chi.Router) {
			r.Get("/", h.SearchAirport)
		})

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetAirport)
			r.Put("/", h.UpdateAirport)
			r.Delete("/", h.DeleteAirport)
		})
	})
}

func (h *AirportHandler) GetAllAirport(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	airports, err := h.service.GetAllAirport(r.Context(), pageSize, offset)
	if err != nil {
		h.logger.Errorf("Failed to get all airports", "error", err)
		respondWithError(w, http.StatusBadRequest, "Failed to get all airports")
		return
	}

	h.logger.Info("All airport data get successfully")
	if airports == nil {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(nil, "No airports found"))
	} else {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(dto.PaginatedResponse{
			Page:     page,
			PageSize: pageSize,
			Data:     airports,
		}, ""))
	}
}

func (h *AirportHandler) GetAirport(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Info("Failed to get airport, invalid id")
		respondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	airport, serviceErr := h.service.GetAirport(r.Context(), id)
	if serviceErr != nil {
		h.logger.Errorf("Failed to get airport", "error", serviceErr)
		respondWithError(w, http.StatusBadRequest, "Failed to get airport")
		return
	}

	h.logger.Info("Airport data get successfully")
	if airport == nil {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(nil, "No airport found"))
	} else {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(airport, ""))
	}
}

func (h *AirportHandler) SearchAirport(w http.ResponseWriter, r *http.Request) {
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

	airports, err := h.service.SearchAirport(r.Context(), icao, facilityName, pageSize, offset)
	if err != nil {
		h.logger.Errorw("Failed to search airports", "error", err)
		respondWithError(w, http.StatusBadRequest, "Failed to search airports")
		return
	}

	h.logger.Info("Airport data searched successfully")
	if airports == nil {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(nil, "No airports found"))
	} else {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(dto.PaginatedResponse{
			Page:     page,
			PageSize: pageSize,
			Data:     airports,
		}, ""))
	}
}

func (h *AirportHandler) CreateAirport(w http.ResponseWriter, r *http.Request) {
	var request dto.Airport

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to create airport, invalid request body")
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.validator.Validate(&request); err != nil {
		h.logger.Errorw("Failed to validate create airport request", "error", err)
		respondWithError(w, http.StatusBadRequest, "Failed to validate create airport request")
		return
	}

	if h.validator.IsComplete(&request) {
		request.Status = "DONE"
	} else {
		request.Status = "PENDING"
	}

	airport, serviceErr := h.service.CreateAirport(r.Context(), &request)
	if serviceErr != nil {
		h.logger.Errorw("Failed to create airport", "error", serviceErr)
		respondWithError(w, http.StatusBadRequest, "Failed to create airport")
		return
	}

	h.logger.Info("Airport data created successfully")
	respondWithJSON(w, http.StatusCreated, dto.NewSuccessResponse(airport, ""))
}

func (h *AirportHandler) UpdateAirport(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Failed to update airport, invalid id")
		respondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	var request dto.Airport
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to update airport, invalid request body")
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()
	request.ID = id

	if err := h.validator.Validate(&request); err != nil {
		h.logger.Errorw("Failed to validate update airport request", "error", err)
		respondWithError(w, http.StatusBadRequest, "Failed to validate update airport request")
		return
	}

	if h.validator.IsComplete(&request) {
		request.Status = "DONE"
	} else {
		request.Status = "PENDING"
	}

	airport, serviceErr := h.service.UpdateAirport(r.Context(), &request)
	if serviceErr != nil {
		h.logger.Errorw("Failed to update airport", "error", serviceErr)
		respondWithError(w, http.StatusBadRequest, "Failed to update airport")
		return
	}

	h.logger.Info("Airport data updated successfully")
	respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(airport, ""))
}

func (h *AirportHandler) DeleteAirport(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Failed to delete airport, invalid id")
		respondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	serviceErr := h.service.DeleteAirport(r.Context(), id)
	if serviceErr != nil {
		h.logger.Errorw("Failed to delete airport", "error", serviceErr)
		respondWithError(w, http.StatusBadRequest, "Failed to delete airport")
		return
	}

	h.logger.Info("Airport data deleted successfully")
	respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(nil, "Airport data deleted successfully"))
}
