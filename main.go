package main

import (
	"context"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config"
	rcpostgres "github.com/andreyloginov-afk/catalog-service/internal/app/repository/conn/postgres"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}

	oldVer, newVer, _, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	if oldVer != newVer {
		log.Info().
			Int64("old_version", oldVer).
			Int64("new_version", newVer).
			Msg("Database migrated")
	} else {
		log.Info().
			Int64("old_version", oldVer).
			Int64("new_version", newVer).
			Msg("Database is up to date")
	}

}
