package category

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type categoryService struct {
	ctxTimeout time.Duration
	repo       repo.Category
}

func NewCategoryService(timeout time.Duration, repository repo.Category) Category {
	return &categoryService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (c *categoryService) beforeCreate(category *entity.Category) {
	category.CreatedAt = time.Now().UTC()
	category.UpdatedAt = time.Now().UTC()
}

func (c *categoryService) beforeUpdate(category *entity.Category) {
	category.UpdatedAt = time.Now().UTC()
}

func (c *categoryService) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	c.beforeCreate(category)

	return c.repo.Create(ctx, category)
}

func (c *categoryService) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	c.beforeUpdate(category)

	return c.repo.Update(ctx, category)
}

func (c *categoryService) Delete(ctx context.Context, categoryID string) error {
	return c.repo.Delete(ctx, categoryID)
}

func (c *categoryService) Get(ctx context.Context, categoryID string) (*entity.Category, error) {
	return c.repo.Get(ctx, categoryID)
}

func (c *categoryService) List(ctx context.Context, page, limit uint64) (*entity.ListCategory, error) {
	return c.repo.List(ctx, page, limit)
}
