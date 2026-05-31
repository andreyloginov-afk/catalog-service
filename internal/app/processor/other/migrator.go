package pprocessor

import (
	"context"
	"sync"

	"github.com/andreyloginov-afk/catalog-service/internal/app/processor"
	"github.com/andreyloginov-afk/catalog-service/internal/app/repository"
	"github.com/rs/zerolog/log"
)

type procMigrate struct {
	migrator repository.Migrate
}

func NewMigrator(migrator repository.Migrate) processor.Processor {
	return &procMigrate{
		migrator: migrator,
	}
}

func (p *procMigrate) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	processor.Wrap(ctx, wg, p.job)
}

func (p *procMigrate) job(ctx context.Context) {
	oldVer, newVer, err := p.migrator.Migrate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("migration failed")
		return
	}

	if oldVer != newVer {
		log.Info().
			Int64("old_version", oldVer).
			Int64("new_version", newVer).
			Msg("database schema updated")
	} else {
		log.Info().Msg("database schema is up to date")
	}
}
