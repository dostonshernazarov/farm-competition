package repo

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Product interface {
	Create(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, productID string) error
	Get(ctx context.Context, params map[string]string) (*entity.Product, error)
	List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListProducts, error)
	UniqueProductName(ctx context.Context, productName string) (int, error)
}
