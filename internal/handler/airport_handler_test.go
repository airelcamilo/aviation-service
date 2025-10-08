package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"aviation-service/internal/dto"
	. "aviation-service/internal/handler"
	. "aviation-service/internal/mock"
	"aviation-service/internal/service"
	utils "aviation-service/internal/testutils"
	"aviation-service/pkg/logger"

	"github.com/go-chi/chi/v5"
)

type mockAirportValidator struct {
	isComplete  bool
	validateErr error
}

func (v *mockAirportValidator) IsComplete(req *dto.Airport) bool {
	return v.isComplete
}

func (v *mockAirportValidator) Validate(req *dto.Airport) error {
	return v.validateErr
}

func TestAirportHandler_GetAllAirport(t *testing.T) {
	tests := []struct {
		name        string
		service     service.IAirportService
		queryParams string
		utils.ExpectedResult
	}{
		{
			name: "Success with data and correct pagination",
			service: &IAirportServiceMock{
				GetAllAirportFunc: func(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KAVL"}}, nil
				},
			},
			queryParams: "?page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data: dto.PaginatedResponse{
					Page:     1,
					PageSize: 10,
					Data:     []dto.Airport{{ID: 1, ICAO: "KAVL"}},
				},
			},
		},
		{
			name: "Success with data and incorrect pagination",
			service: &IAirportServiceMock{
				GetAllAirportFunc: func(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KAVL"}}, nil
				},
			},
			queryParams: "?page=-1&pageSize=A",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data: dto.PaginatedResponse{
					Page:     1,
					PageSize: 10,
					Data:     []dto.Airport{{ID: 1, ICAO: "KAVL"}},
				},
			},
		},
		{
			name: "No data",
			service: &IAirportServiceMock{
				GetAllAirportFunc: func(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
					return nil, nil
				},
			},
			queryParams: "?page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status:  http.StatusOK,
				Message: "No airports found",
			},
		},
		{
			name: "Service error",
			service: &IAirportServiceMock{
				GetAllAirportFunc: func(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{}, fmt.Errorf("DB error")
				},
			},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to get all airports",
			},
		},
	}

	mockAirportValidator := &mockAirportValidator{}
	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAirportHandler(log, tt.service, mockAirportValidator)
			h.RegisterRoutes(r)

			req := httptest.NewRequest(http.MethodGet, "/airport/", nil)
			rr := httptest.NewRecorder()

			h.GetAllAirport(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAirportHandler_GetAirport(t *testing.T) {
	tests := []struct {
		name    string
		service service.IAirportService
		params  map[string]string
		utils.ExpectedResult
	}{
		{
			name: "Success with data",
			service: &IAirportServiceMock{
				GetAirportFunc: func(ctx context.Context, id int) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			params: map[string]string{"id": "1"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data:   &dto.Airport{ID: 1, ICAO: "KAVL"},
			},
		},
		{
			name: "No data",
			service: &IAirportServiceMock{
				GetAirportFunc: func(ctx context.Context, id int) (*dto.Airport, error) {
					return nil, nil
				},
			},
			params: map[string]string{"id": "1"},
			ExpectedResult: utils.ExpectedResult{
				Status:  http.StatusOK,
				Message: "No airport found",
			},
		},
		{
			name: "Service error",
			service: &IAirportServiceMock{
				GetAirportFunc: func(ctx context.Context, id int) (*dto.Airport, error) {
					return nil, fmt.Errorf("DB error")
				},
			},
			params: map[string]string{"id": "1"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to get airport",
			},
		},
		{
			name: "Invalid id",
			service: &IAirportServiceMock{
				GetAirportFunc: func(ctx context.Context, id int) (*dto.Airport, error) {
					return nil, fmt.Errorf("Invalid id")
				},
			},
			params: map[string]string{"id": "A"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Invalid id",
			},
		},
	}

	mockAirportValidator := &mockAirportValidator{}
	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAirportHandler(log, tt.service, mockAirportValidator)
			h.RegisterRoutes(r)

			req := httptest.NewRequest(http.MethodGet, "/airport/", nil)
			req.SetPathValue("id", tt.params["id"])
			rr := httptest.NewRecorder()

			h.GetAirport(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAirportHandler_SearchAirport(t *testing.T) {
	tests := []struct {
		name        string
		service     service.IAirportService
		queryParams string
		utils.ExpectedResult
	}{
		{
			name: "Success with data and correct pagination",
			service: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KAVL"}}, nil
				},
			},
			queryParams: "?icao=KAVL&page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data: dto.PaginatedResponse{
					Page:     1,
					PageSize: 10,
					Data:     []dto.Airport{{ID: 1, ICAO: "KAVL"}},
				},
			},
		},
		{
			name: "Success with data and incorrect pagination",
			service: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return []dto.Airport{{ID: 1, ICAO: "KAVL"}}, nil
				},
			},
			queryParams: "?icao=KAVL&page=-1&pageSize=-1",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data: dto.PaginatedResponse{
					Page:     1,
					PageSize: 10,
					Data:     []dto.Airport{{ID: 1, ICAO: "KAVL"}},
				},
			},
		},
		{
			name: "No data",
			service: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return nil, nil
				},
			},
			queryParams: "?icao=KAVL&page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status:  http.StatusOK,
				Message: "No airports found",
			},
		},
		{
			name: "Service error",
			service: &IAirportServiceMock{
				SearchAirportFunc: func(ctx context.Context, icao, facilityName string, limit, offset int) ([]dto.Airport, error) {
					return nil, fmt.Errorf("DB error")
				},
			},
			queryParams: "?icao=KAVL&page=1&pageSize=10",
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to search airports",
			},
		},
	}

	mockAirportValidator := &mockAirportValidator{}
	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAirportHandler(log, tt.service, mockAirportValidator)
			h.RegisterRoutes(r)

			req := httptest.NewRequest(http.MethodGet, "/airport"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			h.SearchAirport(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAirportHandler_CreateAirport(t *testing.T) {
	tests := []struct {
		name      string
		service   service.IAirportService
		validator AirportValidator
		body      interface{}
		utils.ExpectedResult
	}{
		{
			name: "Valid request Done",
			service: &IAirportServiceMock{
				CreateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			validator: &mockAirportValidator{isComplete: true, validateErr: nil},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusCreated,
				Data:   &dto.Airport{ID: 1, ICAO: "KAVL"},
			},
		},
		{
			name: "Valid request Pending",
			service: &IAirportServiceMock{
				CreateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			validator: &mockAirportValidator{isComplete: false, validateErr: nil},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusCreated,
				Data:   &dto.Airport{ID: 1, ICAO: "KAVL"},
			},
		},
		{
			name: "Invalid request body",
			service: &IAirportServiceMock{
				CreateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			validator: &mockAirportValidator{isComplete: true, validateErr: nil},
			body:      `{"invalid":`,
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Invalid request body",
			},
		},
		{
			name: "Validation failed",
			service: &IAirportServiceMock{
				CreateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			validator: &mockAirportValidator{isComplete: true, validateErr: fmt.Errorf("ICAO is required")},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to validate create airport request",
			},
		},
		{
			name: "Service error",
			service: &IAirportServiceMock{
				CreateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return nil, fmt.Errorf("DB error")
				},
			},
			validator: &mockAirportValidator{isComplete: true, validateErr: nil},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to create airport",
			},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAirportHandler(log, tt.service, tt.validator)
			h.RegisterRoutes(r)

			data, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/airport/", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.CreateAirport(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAirportHandler_UpdateAirport(t *testing.T) {
	tests := []struct {
		name      string
		service   service.IAirportService
		params    map[string]string
		validator AirportValidator
		body      interface{}
		utils.ExpectedResult
	}{
		{
			name: "Valid request Done",
			service: &IAirportServiceMock{
				UpdateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			params:    map[string]string{"id": "1"},
			validator: &mockAirportValidator{isComplete: true, validateErr: nil},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data:   &dto.Airport{ID: 1, ICAO: "KAVL"},
			},
		},
		{
			name: "Valid request Pending",
			service: &IAirportServiceMock{
				UpdateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			params:    map[string]string{"id": "1"},
			validator: &mockAirportValidator{isComplete: false, validateErr: nil},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
				Data:   &dto.Airport{ID: 1, ICAO: "KAVL"},
			},
		},
		{
			name: "Invalid request body",
			service: &IAirportServiceMock{
				UpdateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			params:    map[string]string{"id": "1"},
			validator: &mockAirportValidator{isComplete: true, validateErr: nil},
			body:      `{"invalid":`,
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Invalid request body",
			},
		},
		{
			name: "Validation failed",
			service: &IAirportServiceMock{
				UpdateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return &dto.Airport{ID: 1, ICAO: "KAVL"}, nil
				},
			},
			params:    map[string]string{"id": "1"},
			validator: &mockAirportValidator{isComplete: true, validateErr: fmt.Errorf("ICAO is required")},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to validate update airport request",
			},
		},
		{
			name: "Service error",
			service: &IAirportServiceMock{
				UpdateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return nil, fmt.Errorf("DB error")
				},
			},
			params:    map[string]string{"id": "1"},
			validator: &mockAirportValidator{isComplete: true, validateErr: nil},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to update airport",
			},
		},
		{
			name: "Invalid id",
			service: &IAirportServiceMock{
				UpdateAirportFunc: func(ctx context.Context, req *dto.Airport) (*dto.Airport, error) {
					return nil, fmt.Errorf("Invalid id")
				},
			},
			params:    map[string]string{"id": "A"},
			validator: &mockAirportValidator{isComplete: true, validateErr: nil},
			body:      dto.Airport{ICAO: "KAVL"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Invalid id",
			},
		},
	}

	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAirportHandler(log, tt.service, tt.validator)
			h.RegisterRoutes(r)

			data, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/airport/", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tt.params["id"])
			rr := httptest.NewRecorder()

			h.UpdateAirport(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAirportHandler_DeleteAirport(t *testing.T) {
	tests := []struct {
		name    string
		service service.IAirportService
		params  map[string]string
		utils.ExpectedResult
	}{
		{
			name: "Valid request",
			service: &IAirportServiceMock{
				DeleteAirportFunc: func(ctx context.Context, id int) error {
					return nil
				},
			},
			params: map[string]string{"id": "1"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusOK,
			},
		},
		{
			name: "Service error",
			service: &IAirportServiceMock{
				DeleteAirportFunc: func(ctx context.Context, id int) error {
					return fmt.Errorf("DB error")
				},
			},
			params: map[string]string{"id": "1"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Failed to delete airport",
			},
		},
		{
			name: "Invalid id",
			service: &IAirportServiceMock{
				DeleteAirportFunc: func(ctx context.Context, id int) error {
					return fmt.Errorf("Invalid id")
				},
			},
			params: map[string]string{"id": "A"},
			ExpectedResult: utils.ExpectedResult{
				Status: http.StatusBadRequest,
				Error:  "Invalid id",
			},
		},
	}

	mockAirportValidator := &mockAirportValidator{}
	log := logger.GetLogger()
	defer log.Sync()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			h := NewAirportHandler(log, tt.service, mockAirportValidator)
			h.RegisterRoutes(r)

			req := httptest.NewRequest(http.MethodDelete, "/airport/", nil)
			req.SetPathValue("id", tt.params["id"])
			rr := httptest.NewRecorder()

			h.DeleteAirport(rr, req)

			if err := utils.AssertHandlerResponse(t, rr, tt.ExpectedResult); err != nil {
				t.Error(err)
			}
		})
	}
}
