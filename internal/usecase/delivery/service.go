package delivery

import (
	"context"
	"github.com/google/uuid"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type deliveryService struct {
	ctxTimeout time.Duration
	repo       repo.Delivery
}

func NewDeliveryService(timeout time.Duration, repository repo.Delivery) Delivery {
	return &deliveryService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (a *deliveryService) beforeCreate(delivery *entity.Delivery) {
	delivery.ID = uuid.New().String()
	delivery.CreatedAt = time.Now().UTC()
	delivery.UpdatedAt = time.Now().UTC()
}

func (a *deliveryService) beforeUpdate(delivery *entity.Delivery) {
	delivery.UpdatedAt = time.Now().UTC()
}

func (a *deliveryService) Create(ctx context.Context, delivery *entity.Delivery) (*entity.Delivery, error) {
	a.beforeCreate(delivery)

	return a.repo.Create(ctx, delivery)
}

func (a *deliveryService) Update(ctx context.Context, delivery *entity.Delivery) (*entity.Delivery, error) {
	a.beforeUpdate(delivery)

	return a.repo.Update(ctx, delivery)
}

func (a *deliveryService) Delete(ctx context.Context, deliveryID string) error {
	return a.repo.Delete(ctx, deliveryID)
}

func (a *deliveryService) Get(ctx context.Context, deliveryID string) (*entity.Delivery, error) {
	return a.repo.Get(ctx, deliveryID)
}

func (a *deliveryService) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListDelivery, error) {
	return a.repo.List(ctx, page, limit, params)
}
