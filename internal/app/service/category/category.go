package scategory

import (
	"context"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/andreyloginov-afk/catalog-service/internal/app/repository"
	"github.com/andreyloginov-afk/catalog-service/internal/app/service"
	"github.com/gofrs/uuid"
)

type svc struct {
	repoCategory repository.Category
	repoProduct  repository.Product
}

func NewService(repoCategory repository.Category, repoProduct repository.Product) service.Category {
	return &svc{
		repoCategory: repoCategory,
		repoProduct:  repoProduct,
	}
}

func (s *svc) Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.Category, error) {
	var category entity.Category

	err := s.repoCategory.InsideTx(ctx, func(txCtx context.Context) error {
		existing, err := s.repoCategory.List(txCtx, &req.Name)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return entity.ErrAlreadyExists
		}

		now := time.Now()
		category = entity.Category{
			GUID:      uuid.Must(uuid.NewV4()),
			Name:      req.Name,
			CreatedAt: now,
			UpdatedAt: now,
		}

		return s.repoCategory.Create(txCtx, category)
	})
	if err != nil {
		return entity.Category{}, err
	}
	return category, nil
}

func (s *svc) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error) {
	return s.repoCategory.GetByGUIDs(ctx, guids)
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error) {
	var category entity.Category

	err := s.repoCategory.InsideTx(ctx, func(txCtx context.Context) error {
		categories, err := s.repoCategory.GetByGUIDs(txCtx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(categories) == 0 {
			return entity.ErrNotFound
		}
		category = categories[0]

		existing, err := s.repoCategory.List(txCtx, &req.Name)
		if err != nil {
			return err
		}
		for _, c := range existing {
			if c.GUID != guid {
				return entity.ErrAlreadyExists
			}
		}

		category.Name = req.Name
		category.UpdatedAt = time.Now()

		return s.repoCategory.Update(txCtx, category)
	})
	if err != nil {
		return entity.Category{}, err
	}
	return category, nil
}

func (s *svc) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoCategory.InsideTx(ctx, func(txCtx context.Context) error {
		categories, err := s.repoCategory.GetByGUIDs(txCtx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(categories) == 0 {
			return entity.ErrNotFound
		}

		products, err := s.repoProduct.List(txCtx, nil, &guid, nil, nil)
		if err != nil {
			return err
		}
		if len(products) > 0 {
			return entity.ErrCategoryHasProducts
		}

		return s.repoCategory.Delete(txCtx, guid)
	})
	return err
}

func (s *svc) List(ctx context.Context) ([]entity.Category, error) {
	return s.repoCategory.List(ctx, nil)
}
