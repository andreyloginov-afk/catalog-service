package sproduct

import (
	"context"
	"fmt"
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
	existing, err := s.repoProduct.List(ctx, &req.Name, nil)
	if err != nil {
		return entity.Product{}, err
	}

	if len(existing) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}

	_, err = s.repoCategory.GetByGUID(ctx, req.CategoryGuid)
	if err != nil {
		return entity.Product{}, err
	}

	now := time.Now()
	product := entity.Product{
		GUID:         uuid.Must(uuid.NewV4()),
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGuid: req.CategoryGuid,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.repoProduct.Create(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *svc) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	return s.repoProduct.GetByGUID(ctx, guid)
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	existing, err := s.repoProduct.List(ctx, &req.Name, nil)
	if err != nil {
		return entity.Product{}, err
	}

	if len(existing) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.CategoryGuid = req.CategoryGUID
	product.UpdatedAt = time.Now()

	if err := s.repoProduct.Update(ctx, product); err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *svc) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	products, err := s.repoProduct.List(ctx, nil, &guid)
	if err != nil {
		return fmt.Errorf("failed to list products: %w", err)
	}
	if len(products) > 0 {
		return entity.ErrCategoryHasProducts
	}

	return s.repoProduct.Delete(ctx, guid)

}

func (s *svc) List(ctx context.Context) ([]entity.Product, error) {
	return s.repoProduct.List(ctx, nil, nil)
}
