package handler

import (
	"aviation-service/internal/dto"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AviationSyncHandler struct {
	logger  *zap.SugaredLogger
	service AviationSyncService
}

type AviationSyncService interface {
	Sync(ctx context.Context) (*dto.SyncResponse, error)
}

func NewAviationSyncHandler(logger *zap.SugaredLogger, service AviationSyncService) *AviationSyncHandler {
	return &AviationSyncHandler{
		logger:  logger,
		service: service,
	}
}

func (h *AviationSyncHandler) RegisterRoutes(r chi.Router) {
	r.Route("/sync", func(r chi.Router) {
		r.Post("/", h.Sync)
	})
}

func (h *AviationSyncHandler) Sync(w http.ResponseWriter, r *http.Request) {
	syncResponse, err := h.service.Sync(r.Context())
	if err != nil {
		h.logger.Errorw("Failed to sync", "error", err)
		respondWithError(w, http.StatusBadRequest, "Failed to sync")
		return
	}

	h.logger.Info("Sync airport data successfully")
	if syncResponse == nil {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(nil, "No airport data to sync"))
	} else {
		respondWithJSON(w, http.StatusOK, dto.NewSuccessResponse(syncResponse, ""))
	}
}
