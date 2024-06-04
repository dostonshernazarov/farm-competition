package products

import (
	"context"
	"github.com/google/uuid"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type productService struct {
	ctxTimeout time.Duration
	repo       repo.Product
}

func NewFoodService(timeout time.Duration, repository repo.Product) Product {
	return &productService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (p *productService) beforeCreate(product *entity.Product) {
	product.ID = uuid.New().String()
	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()
}

func (p *productService) beforeUpdate(product *entity.Product) {
	product.UpdatedAt = time.Now().UTC()
}

func (p *productService) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	p.beforeCreate(product)

	return p.repo.Create(ctx, product)
}

func (p *productService) Update(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	p.beforeUpdate(product)

	return p.repo.Update(ctx, product)
}

func (p *productService) Delete(ctx context.Context, productID string) error {
	return p.repo.Delete(ctx, productID)
}

func (p *productService) Get(ctx context.Context, params map[string]string) (*entity.Product, error) {
	return p.repo.Get(ctx, params)
}

func (p *productService) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListProducts, error) {
	return p.repo.List(ctx, page, limit, params)
}

func (p *productService) UniqueProductName(ctx context.Context, productName string) (int, error) {
	return p.repo.UniqueProductName(ctx, productName)
}
