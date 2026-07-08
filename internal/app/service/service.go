package service

import (
	"context"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/gofrs/uuid"
)

type (
	Category interface {
		Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.Category, error)
		GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error)
		Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error)
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context) ([]entity.Category, error)
	}
)

type (
	Product interface {
		Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error)
		GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error)
		Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error)
		Delete(ctx context.Context, guid uuid.UUID) error
		GetByGUIDs(ctx context.Context, guid []uuid.UUID) ([]entity.Product, error)
		List(ctx context.Context) ([]entity.Product, error)
	}
)
