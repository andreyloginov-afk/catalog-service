package main

import (
	"github.com/andreyloginov-afk/catalog-service/internal/app/config"
	rhealth "github.com/andreyloginov-afk/catalog-service/internal/app/handler/health"
	rprocessor "github.com/andreyloginov-afk/catalog-service/internal/app/processor/http"
	"github.com/rs/zerolog/log"
)

func main() {
	// Загружаем конфигурацию
	config.Load()
	cfg := config.Root

	// Создаём handler health-check
	healthHandler := rhealth.NewHandler()

	// Создаём HTTP сервер
	httpServer := rprocessor.NewHttp(
		healthHandler,
		cfg.Processor.WebServer,
	)

	// Запускаем сервер
	if err := httpServer.Serve(); err != nil {
		log.Fatal().
			Err(err).
			Msg("HTTP server failed")
	}
}
