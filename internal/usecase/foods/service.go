package foods

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type foodService struct {
	ctxTimeout time.Duration
	repo       repo.Food
}

func NewFoodService(timeout time.Duration, repository repo.Food) Food {
	return &foodService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (f *foodService) beforeCreate(food *entity.Food) {
	food.CreatedAt = time.Now().UTC()
	food.UpdatedAt = time.Now().UTC()
}

func (f *foodService) beforeUpdate(food *entity.Food) {
	food.UpdatedAt = time.Now().UTC()
}

func (f *foodService) Create(ctx context.Context, food *entity.Food) (*entity.Food, error) {
	f.beforeCreate(food)

	return f.repo.Create(ctx, food)
}

func (f *foodService) Update(ctx context.Context, food *entity.Food) (*entity.Food, error) {
	f.beforeUpdate(food)

	return f.repo.Update(ctx, food)
}

func (f *foodService) Delete(ctx context.Context, foodID string) error {
	return f.repo.Delete(ctx, foodID)
}

func (f *foodService) Get(ctx context.Context, foodID string) (*entity.Food, error) {
	return f.repo.Get(ctx, foodID)
}

func (f *foodService) List(ctx context.Context, page, limit uint64) (*entity.ListFoods, error) {
	return f.repo.List(ctx, page, limit)
}
