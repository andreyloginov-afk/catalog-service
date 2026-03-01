package config

import (
	"fmt"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Repository section.Repository `split_words:"true"`
	Processor  section.Processor  `split_words:"true"`
	Monitor    section.Monitor    `split_words:"true"`
}

func Load() (*Config, error) {
	// Загружаем .env
	_ = godotenv.Load()

	var cfg Config

	if err := envconfig.Process("APP", &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
