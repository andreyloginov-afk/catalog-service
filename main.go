package main

import (
	"context"
	"os"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config"
	hcategory "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http/category"
	rhealth "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http/health"
	hproduct "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http/product"
	rprocessor "github.com/andreyloginov-afk/catalog-service/internal/app/processor/http"
	pcategory "github.com/andreyloginov-afk/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/andreyloginov-afk/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/andreyloginov-afk/catalog-service/internal/app/repository/product"
	scategory "github.com/andreyloginov-afk/catalog-service/internal/app/service/category"
	sproduct "github.com/andreyloginov-afk/catalog-service/internal/app/service/product"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	//конфиг
	config.Load(config.LoadArgs{
		Output:          os.Stdout,
		EnableSimpleLog: true,
	})
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
	// Репозитории
	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	// Сервисы
	categorySvc := scategory.NewService(categoryRepo, productRepo)
	productSvc := sproduct.NewService(productRepo, categoryRepo)

	// Хендлеры
	healthHandler := rhealth.NewHandler()
	categoryHandler := hcategory.NewHandler(categorySvc)
	productHandler := hproduct.NewHandler(productSvc)

	// сервер
	server := rprocessor.NewHttp(
		healthHandler,
		categoryHandler,
		productHandler,
		cfg.Processor.WebServer,
	)

	if err := server.Serve(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
