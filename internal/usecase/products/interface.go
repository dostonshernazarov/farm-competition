package products

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Product interface {
	Create(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, productID string) error
	Get(ctx context.Context, productID string) (*entity.Product, error)
	List(ctx context.Context, page, limit uint64) (*entity.ListProducts, error)
}
