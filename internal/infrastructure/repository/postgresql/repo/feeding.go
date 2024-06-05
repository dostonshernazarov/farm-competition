package repo

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Feeding interface {
	Create(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error)
	Update(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error)
	Delete(ctx context.Context, feedingID string) error
}
