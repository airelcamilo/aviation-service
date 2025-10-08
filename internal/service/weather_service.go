package service

import (
	"aviation-service/config"
	"aviation-service/internal/dto"
	"aviation-service/internal/utils"
	"aviation-service/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"go.uber.org/zap"
)

//go:generate moq -out ../mock/weather_service_mock.go -pkg=mock . IWeatherService
type IWeatherService interface {
	GetWeather(ctx context.Context, city string) (*dto.Weather, error)
}

type WeatherService struct {
	logger      *zap.SugaredLogger
	cfg         config.Config
	client      Client
	redisClient redis.RedisClient
}

func NewWeatherService(logger *zap.SugaredLogger, cfg config.Config, client Client, redisClient redis.RedisClient) *WeatherService {
	return &WeatherService{
		logger:      logger,
		cfg:         cfg,
		client:      client,
		redisClient: redisClient,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (*dto.Weather, error) {
	cacheKey := fmt.Sprintf("weather:%s", city)
	var weather dto.WeatherDataResponse
	s.logger.Infow("Weather cache hit", "city", city)
	cacheErr := utils.GetStruct(s.redisClient, ctx, cacheKey, &weather)
	if cacheErr == nil {
		return &weather.Current, nil
	}
	s.logger.Infow("No airport data from cache, fetching from API", "error", cacheErr)

	s.logger.Infow("Fetching weather data", "city", city)
	params := url.Values{}
	params.Add("key", s.cfg.WEATHER_API_KEY)
	params.Add("q", city)
	resp, err := s.client.Get(s.cfg.WEATHER_API_URL + "/current.json?" + params.Encode())
	if err != nil {
		s.logger.Errorw("Error fetching weather data", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		s.logger.Errorw("Error reading response body", "error", readErr)
		return nil, readErr
	}

	jsonErr := json.Unmarshal(bodyBytes, &weather)
	if jsonErr != nil {
		s.logger.Errorw("Error decoding body", "error", jsonErr)
		return nil, jsonErr
	}

	exp := calculateTTL()
	if err := s.redisClient.Set(ctx, cacheKey, string(bodyBytes), exp).Err(); err != nil {
		s.logger.Errorw("Error to cache weather", "error", err)
	}
	return &weather.Current, nil
}

func calculateTTL() time.Duration {
	now := time.Now()
	min := now.Minute()

	nextQuarter := ((min/15)+1)*15 - min
	return time.Duration(nextQuarter) * time.Minute
}
