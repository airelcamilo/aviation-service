package main

import (
    "context"
    "time"
    "net/http"

    "github.com/robfig/cron/v3"
    "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"aviation-service/config"
	"aviation-service/internal/repository"
	"aviation-service/internal/service"
	
	"aviation-service/pkg/logger"
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

    airportRepo := repository.NewAirportRepository(db)
    httpClient := http.DefaultClient
    airportService := service.NewAirportService(log, airportRepo, cfg, httpClient, nil)
    aviationSyncService := service.NewAviationSyncService(log, airportRepo, airportService)

    c := cron.New()
    c.AddFunc("0 5 * * *", func() {
        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
        defer cancel()

        syncResponse, err := aviationSyncService.Sync(ctx)
        if err != nil {
            log.Errorw("Failed to sync", "error", err)
            return
        }

        if syncResponse == nil {
            log.Info("No airport data to sync")
        } else {
            log.Infow("Sync airport data successfully", "sync", syncResponse)
        }
    })

    log.Info("Starting aviation sync cron scheduler")
    c.Start()
    select {}
}
