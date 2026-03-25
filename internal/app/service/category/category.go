package scategory

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
	// Бизнес валидация: проверка дубликата имени
	existing, err := s.repoCategory.List(ctx, &req.Name)
	if err != nil {
		return entity.Category{}, err
	}
	if len(existing) > 0 {
		return entity.Category{}, entity.ErrAlreadyExists
	}
	//подготовка entity
	now := time.Now()
	category := entity.Category{
		GUID:      uuid.Must(uuid.NewV4()),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	// сохранение
	if err := s.repoCategory.Create(ctx, category); err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *svc) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	return s.repoCategory.GetByGUID(ctx, guid)
	// репозиторий уже возвращает ErrNotFound
}

func (s *svc) Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error) {
	// получает существующую категорию через GetByGUID
	category, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Category{}, err
	}
	// проверкам на дубликаты через List
	existing, err := s.repoCategory.List(ctx, &req.Name)
	if err != nil {
		return entity.Category{}, err
	}
	for _, c := range existing {
		if c.GUID != guid {
			return entity.Category{}, entity.ErrAlreadyExists
		}
	}
	// обновляет поля
	category.Name = req.Name
	category.UpdatedAt = time.Now()

	if err := s.repoCategory.Update(ctx, category); err != nil {
		return entity.Category{}, err
	}

	return category, nil

}

func (s *svc) Delete(ctx context.Context, guid uuid.UUID) error {
	// проверяет существования
	_, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}
	// проверяет наличие
	products, err := s.repoProduct.List(ctx, nil, &guid)
	if err != nil {
		return fmt.Errorf("failed to list products: %w", err)
	}
	if len(products) > 0 {
		return entity.ErrCategoryHasProducts
	}

	return s.repoCategory.Delete(ctx, guid)
}

func (s *svc) List(ctx context.Context) ([]entity.Category, error) {
	return s.repoCategory.List(ctx, nil)
}
