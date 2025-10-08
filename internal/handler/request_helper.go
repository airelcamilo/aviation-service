package handler

import (
	"encoding/json"
	"net/http"

	"aviation-service/internal/dto"
	"aviation-service/pkg/logger"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Errorw("Error marshal", "error", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		logger.Errorw("Error write", "error", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message interface{}) {
	respondWithJSON(w, code, dto.NewErrorResponse(message, "Error occurred"))
}
