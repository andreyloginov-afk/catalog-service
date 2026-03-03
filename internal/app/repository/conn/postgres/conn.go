package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
	"github.com/andreyloginov-afk/catalog-service/migration"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB
		cfg      section.RepositoryPostgres
	}
	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func (c *Client) Close() error {
	if c.rawBunDB != nil {
		return c.rawBunDB.Close()
	}
	return nil
}

func NewConn(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	u := &url.URL{
		Scheme: "postgres",
		Host:   cfg.Address,
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Path:   cfg.Name,
	}

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	sqldb.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqldb, pgdialect.New())

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqldb.PingContext(pingCtx); err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Client{
		_bunDB:   bunDB,
		rawBunDB: bunDB,
		cfg:      cfg,
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer, applied int64, err error) {
	migrations := migrate.NewMigrations()
	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, 0, fmt.Errorf("discover migrations: %w", err)
	}
	//тут опция мигратора с таблицой
	migrator := migrate.NewMigrator(
		c.rawBunDB,
		migrations,
		migrate.WithTableName("catalog_migrations"),
	)

	if err = migrator.Init(ctx); err != nil {
		return 0, 0, 0, fmt.Errorf("migrator init: %w", err)
	}

	// Тут версия ДО миграции
	appliedMigrations, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get applied migrations: %w", err)
	}
	oldVer = appliedMigrations.LastGroupID()

	//  тут применяем миграции
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return oldVer, oldVer, 0, fmt.Errorf("migrate: %w", err)
	}

	// Если ничего не применилось
	if group.IsZero() {
		return oldVer, oldVer, 0, nil
	}

	newVer = group.ID
	applied = int64(len(group.Migrations))

	return oldVer, newVer, applied, nil
}
