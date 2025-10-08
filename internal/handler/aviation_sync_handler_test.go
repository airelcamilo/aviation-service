package handler_test

import (
	"aviation-service/internal/dto"
	utils "aviation-service/internal/testutils"
	. "aviation-service/internal/handler"
	"aviation-service/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mockSyncService struct {
	response *dto.SyncResponse
	err      error
}

func (m *mockSyncService) Sync(ctx context.Context) (*dto.SyncResponse, error) {
	return m.response, m.err
}

func TestAviationSyncHandler_Sync(t *testing.T) {
	tests := []struct {
		name    string
		service AviationSyncService
		utils.ExpectedResult
	}{
		{
			name:    "Success with data",
			service: &mockSyncService{response: &dto.SyncResponse{Total: 12, Success: 10, Failed: 2}, err: nil},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data:   dto.SyncResponse{Total: 12, Success: 10, Failed: 2},
			},
		},
		{
			name:    "Success no data",
			service: &mockSyncService{response: nil, err: nil},
			ExpectedResult: utils.ExpectedResult{
				Status:  http.StatusOK,
				Data:    nil,
				Message: "No airport data to sync",
			},
		},
		{
			name:    "Failed to sync",
			service: &mockSyncService{response: nil, err: fmt.Errorf("Sync error")},
			ExpectedResult: utils.ExpectedResult{
				Status:  http.StatusBadRequest,
				Data:    nil,
				Message: "Error occurred",
				Error:   "Failed to sync",
			},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAviationSyncHandler(log, tt.service)
			h.RegisterRoutes(r)

			req := httptest.NewRequest(http.MethodPost, "/sync", bytes.NewBuffer(nil))
			rr := httptest.NewRecorder()

			h.Sync(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}
