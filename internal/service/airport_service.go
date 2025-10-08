package service

import (
	"aviation-service/config"
	"aviation-service/internal/dto"
	"aviation-service/internal/repository"
	"aviation-service/internal/utils"
	"aviation-service/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

//go:generate moq -out ../mock/airport_service_mock.go -pkg=mock . IAirportService
type IAirportService interface {
	GetAllAirport(ctx context.Context, limit, offset int) ([]dto.Airport, error)
	GetAirport(ctx context.Context, id int) (*dto.Airport, error)
	CreateAirport(ctx context.Context, request *dto.Airport) (*dto.Airport, error)
	SearchAirport(ctx context.Context, icao string, name string, limit, offset int) ([]dto.Airport, error)
	UpdateAirport(ctx context.Context, request *dto.Airport) (*dto.Airport, error)
	DeleteAirport(ctx context.Context, id int) error
	FetchAirportData(icaos string) (*dto.AirportDataResponse, error)
}

type Client interface {
	Get(url string) (*http.Response, error)
}

type AirportService struct {
	logger      *zap.SugaredLogger
	airportRepo repository.IAirportRepository
	cfg         config.Config
	client      Client
	redisClient redis.RedisClient
}

func NewAirportService(logger *zap.SugaredLogger, airportRepo repository.IAirportRepository, cfg config.Config, client Client, redisClient redis.RedisClient) *AirportService {
	return &AirportService{
		logger:      logger,
		airportRepo: airportRepo,
		cfg:         cfg,
		client:      client,
		redisClient: redisClient,
	}
}

func (s *AirportService) GetAllAirport(ctx context.Context, limit, offset int) ([]dto.Airport, error) {
	airports, err := s.airportRepo.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Errorw("Failed to get all airports", "error", err)
		return nil, err
	}
	return airports, nil
}

func (s *AirportService) GetAirport(ctx context.Context, id int) (*dto.Airport, error) {
	airport, err := s.airportRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get airport", "error", err)
		return nil, err
	}
	return airport, nil
}

func (s *AirportService) SearchAirport(ctx context.Context, icao string, facilityName string, limit, offset int) ([]dto.Airport, error) {
	cacheKey := fmt.Sprintf("airport:%s:%s:%s:%s", icao, facilityName, limit, offset)
	var airports []dto.Airport
	s.logger.Infow("Airport cache hit", "icao", icao, "facilityName", facilityName)
	cacheErr := utils.GetStruct(s.redisClient, ctx, cacheKey, &airports)
	if cacheErr == nil {
		return airports, nil
	}
	s.logger.Infow("No airport data from cache, fetching from repo", "error", cacheErr)

	s.logger.Infow("Get airports data from repo", "icao", icao, "facilityName", facilityName)
	airports, err := s.airportRepo.GetByICAOOrFacilityName(ctx, icao, facilityName, limit, offset)
	if err != nil {
		s.logger.Errorw("Failed to get airports from repo", "error", err, "icao", icao, "facilityName", facilityName)
		return nil, err
	}

	if len(airports) > 0 || icao == "" {
		if err := utils.SetStruct(s.redisClient, ctx, cacheKey, airports, 24*time.Hour); err != nil {
			s.logger.Infow("Error set cache", "error", err)
		}
		return airports, nil
	}
	s.logger.Infow("No airport data from repo, fetching from API", "error", cacheErr)

	airportResponse, err := s.FetchAirportData(icao)
	if err != nil {
		s.logger.Errorw("Failed to get airports from API", "error", err, "icao", icao)
		return nil, err
	}

	airports = (*airportResponse)[icao]
	if len(airports) == 0 {
		return nil, nil
	}

	airports[0].Status = "DONE"
	inserted, err := s.airportRepo.Insert(ctx, &airports[0])
	if err != nil {
		s.logger.Errorw("Failed to insert airport from API", "error", err, "icao", icao)
		return nil, err
	}
	airports = []dto.Airport{*inserted}

	if err := utils.SetStruct(s.redisClient, ctx, cacheKey, airports, 24*time.Hour); err != nil {
		s.logger.Infow("Error set cache", "error", err)
	}
	return airports, nil
}

func (s *AirportService) CreateAirport(ctx context.Context, request *dto.Airport) (*dto.Airport, error) {
	airport, err := s.airportRepo.Insert(ctx, request)
	if err != nil {
		s.logger.Errorw("Failed to create airport", "error", err)
		return nil, err
	}
	return airport, nil
}

func (s *AirportService) UpdateAirport(ctx context.Context, request *dto.Airport) (*dto.Airport, error) {
	airport, err := s.airportRepo.UpdateById(ctx, request)
	if err != nil {
		s.logger.Errorw("Failed to update airport", "error", err)
		return nil, err
	}
	return airport, nil
}

func (s *AirportService) DeleteAirport(ctx context.Context, id int) error {
	err := s.airportRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to delete airport", "error", err)
		return err
	}
	return nil
}

func (s *AirportService) FetchAirportData(icaos string) (*dto.AirportDataResponse, error) {
	s.logger.Infow("Fetching airports data", "icaos", icaos)
	params := url.Values{}
	params.Add("apt", icaos)
	resp, err := s.client.Get(s.cfg.AIRPORT_API_URL + "/airports?" + params.Encode())
	if err != nil {
		s.logger.Errorw("Error fetching airports data", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	var airports dto.AirportDataResponse
	jsonErr := json.NewDecoder(resp.Body).Decode(&airports)
	if jsonErr != nil {
		s.logger.Errorw("Error decoding body", "error", jsonErr)
		return nil, jsonErr
	}
	return &airports, nil
}
