package main

import (
	"context"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config"
	rcpostgres "github.com/andreyloginov-afk/catalog-service/internal/app/repository/conn/postgres"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	//конфиг
	config.Load()
	cfg := config.Root

	// PostgreSQL
	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to PostgreSQL")
	}
	defer func() {
		if err := pgClient.Close(); err != nil {
			log.Error().
				Err(err).
				Msg("failed to close PostgreSQL connection")
		}
	}()

	// миграции
	oldVer, newVer, applied, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to run migrations")
	}

	if applied > 0 {
		log.Info().
			Int64("old_version", oldVer).
			Int64("new_version", newVer).
			Int64("applied", applied).
			Msg("database migrated")
	} else {
		log.Info().
			Int64("version", oldVer).
			Msg("database is up to date")
	}
}
