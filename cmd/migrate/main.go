package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"aviation-service/config"
	"aviation-service/pkg/seeding"
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

	m, err := migrate.New(
		"file://migrations",
		cfg.DATABASE_URL)
	if err != nil {
		logger.Fatalw("Failed to make migration", "error", err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			logger.Fatalw("Failed to migration up", "error", err)
		}
	}
	logger.Info("Success to migrate")

	if err := seeding.RunSeeding(db, "seedings"); err != nil {
		logger.Fatalw("Failed seeding", "error", err)
	}
	logger.Info("Success to seed")
}
