package foods

import (
	"context"
	"github.com/google/uuid"
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
	food.ID = uuid.New().String()
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

func (f *foodService) Get(ctx context.Context, params map[string]string) (*entity.Food, error) {
	return f.repo.Get(ctx, params)
}

func (f *foodService) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListFoods, error) {
	return f.repo.List(ctx, page, limit, params)
}

func (f *foodService) UniqueFoodName(ctx context.Context, foodName string) (int, error) {
	return f.repo.UniqueFoodName(ctx, foodName)
}
