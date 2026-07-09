package sproduct

import (
	"context"
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
		existing, err := s.repoProduct.List(txCtx, &req.Name, nil, nil, nil)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return entity.ErrAlreadyExists
		}

		categories, err := s.repoCategory.GetByGUIDs(txCtx, []uuid.UUID{req.CategoryGuid})
		if err != nil {
			return err
		}
		if len(categories) == 0 {
			return entity.ErrNotFound
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

func (s *svc) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error) {
	return s.repoProduct.GetByGUIDs(ctx, guids)
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	var product entity.Product

	err := s.repoProduct.InsideTx(ctx, func(txCtx context.Context) error {
		products, err := s.repoProduct.GetByGUIDs(txCtx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(products) == 0 {
			return entity.ErrNotFound
		}
		product = products[0]

		existing, err := s.repoProduct.List(txCtx, &req.Name, nil, nil, nil)
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
		products, err := s.repoProduct.GetByGUIDs(txCtx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(products) == 0 {
			return entity.ErrNotFound
		}

		return s.repoProduct.Delete(txCtx, guid)
	})
}

func (s *svc) List(ctx context.Context, req entity.RequestProductList) ([]entity.Product, error) {
	return s.repoProduct.List(ctx, req.Name, req.CategoryGUID, req.MinPrice, req.MaxPrice)
}
