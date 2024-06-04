package repo

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Delivery interface {
	Create(ctx context.Context, delivery *entity.Delivery) (*entity.Delivery, error)
	Update(ctx context.Context, delivery *entity.Delivery) (*entity.Delivery, error)
	Delete(ctx context.Context, animalID string) error
	Get(ctx context.Context, deliveryID string) (*entity.Delivery, error)
	List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListDelivery, error)
}
