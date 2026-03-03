package section

import "github.com/andreyloginov-afk/catalog-service/internal/app/util"

type (
	Repository struct {
		Postgres RepositoryPostgres
	}
	RepositoryPostgres struct {
		Address        string        `required:"true" default:"localhost" split_words:"true"`
		Username       string        `required:"true" split_words:"true"`
		Password       string        `required:"true" split_words:"true"`
		Name           string        `required:"true" split_words:"true"`
		ReadTimeout    util.Duration `required:"true" split_words:"true"`
		WriteTimeout   util.Duration `required:"true" split_words:"true"`
		MigrationTable string        `split_words:"true" default:"schema_migrations"`
	}
)
