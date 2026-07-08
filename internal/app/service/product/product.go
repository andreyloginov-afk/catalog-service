package sproduct

import (
	"context"
	"errors"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/andreyloginov-afk/catalog-service/internal/app/repository"
	"github.com/andreyloginov-afk/catalog-service/internal/app/service"
	"github.com/gofrs/uuid"
)

type svc struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) service.Product {
	return &svc{
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (s *svc) Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error) {
	var product entity.Product

	err := s.repoProduct.InsideTx(ctx, func(txCtx context.Context) error {
		existing, err := s.repoProduct.List(txCtx, &req.Name, nil)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return entity.ErrAlreadyExists
		}

		_, err = s.repoCategory.GetByGUID(txCtx, req.CategoryGuid)
		if err != nil {
			return err
		}

		now := time.Now()
		product = entity.Product{
			GUID:         uuid.Must(uuid.NewV4()),
			Name:         req.Name,
			Description:  req.Description,
			Price:        req.Price,
			CategoryGuid: req.CategoryGuid,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		return s.repoProduct.Create(txCtx, product)
	})
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *svc) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	return s.repoProduct.GetByGUID(ctx, guid)
}

func (s *svc) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error) {
	products := make([]entity.Product, 0, len(guids))
	for _, guid := range guids {
		p, err := s.repoProduct.GetByGUID(ctx, guid)
		if errors.Is(err, entity.ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	var product entity.Product

	err := s.repoProduct.InsideTx(ctx, func(txCtx context.Context) error {
		var err error
		product, err = s.repoProduct.GetByGUID(txCtx, guid)
		if err != nil {
			return err
		}

		existing, err := s.repoProduct.List(txCtx, &req.Name, nil)
		if err != nil {
			return err
		}
		for _, p := range existing {
			if p.GUID != guid {
				return entity.ErrAlreadyExists
			}
		}

		product.Name = req.Name
		product.Description = req.Description
		product.Price = req.Price
		product.CategoryGuid = req.CategoryGUID
		product.UpdatedAt = time.Now()

		return s.repoProduct.Update(txCtx, product)
	})
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *svc) Delete(ctx context.Context, guid uuid.UUID) error {
	return s.repoProduct.InsideTx(ctx, func(txCtx context.Context) error {
		if _, err := s.repoProduct.GetByGUID(txCtx, guid); err != nil {
			return err
		}

		return s.repoProduct.Delete(txCtx, guid)
	})
}

func (s *svc) List(ctx context.Context) ([]entity.Product, error) {
	return s.repoProduct.List(ctx, nil, nil)
}
