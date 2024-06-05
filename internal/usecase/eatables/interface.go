package eatables

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Eatable interface {
	Create(ctx context.Context, eatable *entity.Eatables) (*entity.EatablesRes, error)
	Update(ctx context.Context, eatable *entity.Eatables) (*entity.EatablesRes, error)
	Delete(ctx context.Context, eatableID string) error
	GetDrugs(ctx context.Context, page, limit uint64, animalID string) (*entity.ListDrugEatables, error)
	GetFoods(ctx context.Context, page, limit uint64, animalID string) (*entity.ListFoodEatables, error)
}
