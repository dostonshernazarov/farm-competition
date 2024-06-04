package foods

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Food interface {
	Create(ctx context.Context, food *entity.Food) (*entity.Food, error)
	Update(ctx context.Context, food *entity.Food) (*entity.Food, error)
	Delete(ctx context.Context, foodID string) error
	Get(ctx context.Context, params map[string]string) (*entity.Food, error)
	List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListFoods, error)
	UniqueFoodName(ctx context.Context, foodName string) (int, error)
}
