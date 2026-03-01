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

type Client struct {
	bunDB *bun.DB
	cfg   section.RepositoryPostgres
	db    *sql.DB
}

func (c *Client) GetRawBunDB() *bun.DB {
	return c.bunDB
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
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
	sqldb.SetConnMaxLifetime(time.Hour)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqldb.PingContext(pingCtx); err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	bunDB := bun.NewDB(sqldb, pgdialect.New())

	return &Client{
		bunDB: bunDB,
		cfg:   cfg,
		db:    sqldb,
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer, applied int64, err error) {
	migrations := migrate.NewMigrations()
	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, 0, fmt.Errorf("discover migrations: %w", err)
	}

	migrator := migrate.NewMigrator(c.bunDB, migrations)

	if err = migrator.Init(ctx); err != nil {
		return 0, 0, 0, fmt.Errorf("migrator init: %w", err)
	}

	if err = migrator.Lock(ctx); err != nil {
		return 0, 0, 0, fmt.Errorf("migrator lock: %w", err)
	}
	defer func() {
		if unlockErr := migrator.Unlock(ctx); unlockErr != nil && err == nil {
			err = fmt.Errorf("migrator unlock: %w", unlockErr)
		}
	}()

	before, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get applied migrations (before): %w", err)
	}
	oldVer = before.LastGroupID()

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return oldVer, oldVer, 0, fmt.Errorf("migrate: %w", err)
	}

	after, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return oldVer, oldVer, 0, fmt.Errorf("get applied migrations (after): %w", err)
	}
	newVer = after.LastGroupID()

	if !group.IsZero() {
		applied = int64(len(group.Migrations))
	}

	return oldVer, newVer, applied, err
}
