package foods

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Food interface {
	Create(ctx context.Context, food *entity.Food) (*entity.Food, error)
	Update(ctx context.Context, food *entity.Food) (*entity.Food, error)
	Delete(ctx context.Context, foodID string) error
	Get(ctx context.Context, foodID string) (*entity.Food, error)
	List(ctx context.Context, page, limit uint64) (*entity.ListFoods, error)
}
