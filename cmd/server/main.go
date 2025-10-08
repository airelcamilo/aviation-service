package main

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"aviation-service/config"
	"aviation-service/internal/handler"
	"aviation-service/internal/repository"
	"aviation-service/internal/service"
	"aviation-service/internal/utils"
	
	"aviation-service/pkg/httpserver"
	"aviation-service/pkg/logger"
	"aviation-service/pkg/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalw("Failed to load configuration", "error", err)
	}

	log := logger.GetLogger()
	defer log.Sync()

	db, err := sqlx.Connect("postgres", cfg.DATABASE_URL)
	if err != nil {
		logger.Fatalw("Failed to connect to db", "error", err)
	}
	defer db.Close()

	redisClient, err := redis.NewRedisClient(cfg.REDIS_URL)
	if err != nil {
		logger.Fatalw("Failed to connect to Redis", "error", err)
	}
	defer redisClient.Close()

	airportRepo := repository.NewAirportRepository(db)
	client := http.DefaultClient

	airportService := service.NewAirportService(log, airportRepo, cfg, client, redisClient)
	aviationSyncService := service.NewAviationSyncService(log, airportRepo, airportService)
	weatherService := service.NewWeatherService(log, cfg, client, redisClient)
	airportWeatherService := service.NewAirportWeatherService(log, airportService, weatherService)

	airportValidator := utils.NewAirportValidator()
	airportHandler := handler.NewAirportHandler(log, airportService, airportValidator)
	aviationSyncHandler := handler.NewAviationSyncHandler(log, aviationSyncService)
	weatherHandler := handler.NewWeatherHandler(log, weatherService)
	airportWeatherHandler := handler.NewAirportWeatherHandler(log, airportWeatherService)

	router := httpserver.NewRouter(
		airportHandler,
		aviationSyncHandler,
		weatherHandler,
		airportWeatherHandler,
	)

	server := httpserver.NewServer(router, "8000")
	logger.Info("Starting HTTP server")
	server.Start()
}
