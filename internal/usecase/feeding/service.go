package feeding

import (
	"context"
	"github.com/google/uuid"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type feedingService struct {
	ctxTimeout time.Duration
	repo       repo.Feeding
}

func NewFeedingService(timeout time.Duration, repository repo.Feeding) Feeding {
	return &feedingService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (f *feedingService) beforeCreate(feeding *entity.Feeding) {
	feeding.ID = uuid.New().String()
	feeding.CreatedAt = time.Now().UTC()
	feeding.UpdatedAt = time.Now().UTC()
}

func (d *feedingService) beforeUpdate(feeding *entity.Feeding) {
	feeding.UpdatedAt = time.Now().UTC()
}

func (d *feedingService) Create(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error) {
	d.beforeCreate(feeding)

	return d.repo.Create(ctx, feeding)
}

func (d *feedingService) Update(ctx context.Context, feeding *entity.Feeding) (*entity.FeedingRes, error) {
	d.beforeUpdate(feeding)

	return d.repo.Update(ctx, feeding)
}

func (d *feedingService) Delete(ctx context.Context, feedingID string) error {
	return d.repo.Delete(ctx, feedingID)
}
