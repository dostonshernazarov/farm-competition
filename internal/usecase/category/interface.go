package category

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Category interface {
	Create(ctx context.Context, category *entity.Category) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) (*entity.Category, error)
	Delete(ctx context.Context, categoryID string) error
	Get(ctx context.Context, categoryID string) (*entity.Category, error)
	List(ctx context.Context, page, limit uint64) (*entity.ListCategory, error)
}
