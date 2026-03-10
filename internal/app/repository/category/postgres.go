package category

import (
	"context"
	"database/sql"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/andreyloginov-afk/catalog-service/internal/app/repository"
	rcpostgres "github.com/andreyloginov-afk/catalog-service/internal/app/repository/conn/postgres"
	"github.com/andreyloginov-afk/catalog-service/internal/app/util"
	"github.com/gofrs/uuid"
)

type (
	repoPg struct {
		*_DB
	}

	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Category {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, category entity.Category) error {
	_, err := r.NewInsert().
		Model(&category).
		Exec(ctx)

	return err
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	var category entity.Category

	err := r.NewSelect().
		Model(&category).
		Where("guid = ?", guid).
		Scan(ctx)

	if err != nil {
		return entity.Category{}, util.ReplaceErr1(err, sql.ErrNoRows, entity.ErrNotFound)
	}

	return category, nil
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) error {
	res, err := r.NewUpdate().
		Model(&category).
		WherePK().
		ExcludeColumn("id", "created_at").
		Exec(ctx)

	return rcpostgres.UpdateErr(res, err)
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r.NewDelete().
		Model((*entity.Category)(nil)).
		Where("guid = ?", guid).
		Exec(ctx)

	if err != nil {
		return rcpostgres.DeleteErr(err)
	}

	return nil
}

func (r *repoPg) List(ctx context.Context, name *string) ([]entity.Category, error) {
	var categories []entity.Category

	query := r.NewSelect().Model(&categories)

	if name != nil {
		query = query.Where("name = ?", *name)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
