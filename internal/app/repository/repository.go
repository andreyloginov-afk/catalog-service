package repository

import (
	"context"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/gofrs/uuid"
)

type Category interface {
	Create(ctx context.Context, category entity.Category) error
	GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error)
	Update(ctx context.Context, category entity.Category) error
	Delete(ctx context.Context, guid uuid.UUID) error
	List(ctx context.Context, name *string) ([]entity.Category, error)
}

type Product interface {
	Create(ctx context.Context, product entity.Product) error
	GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error)
	Update(ctx context.Context, product entity.Product) error
	Delete(ctx context.Context, guid uuid.UUID) error
	List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error)
}
