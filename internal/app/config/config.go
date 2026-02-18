package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository `split_words:"true"`
	Processor  section.Processor  `split_words:"true"`
	Monitor    section.Monitor    `split_words:"true"`
}

var Root Config

func Load() {
	_ = godotenv.Load()

	err := envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatal("Failed to parse config:", err)
	}
}
