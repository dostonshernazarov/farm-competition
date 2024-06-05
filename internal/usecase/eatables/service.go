package eatables

import (
	"context"
	"github.com/google/uuid"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type eatablesService struct {
	ctxTimeout time.Duration
	repo       repo.Eatable
}

func NewEatableService(timeout time.Duration, repository repo.Eatable) Eatable {
	return &eatablesService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (d *eatablesService) beforeCreate(drug *entity.Eatables) {
	drug.ID = uuid.New().String()
	drug.CreatedAt = time.Now().UTC()
	drug.UpdatedAt = time.Now().UTC()
}

func (d *eatablesService) beforeUpdate(drug *entity.Eatables) {
	drug.UpdatedAt = time.Now().UTC()
}

func (d *eatablesService) Create(ctx context.Context, drug *entity.Eatables) (*entity.EatablesRes, error) {
	d.beforeCreate(drug)

	return d.repo.Create(ctx, drug)
}

func (d *eatablesService) Update(ctx context.Context, drug *entity.Eatables) (*entity.EatablesRes, error) {
	d.beforeUpdate(drug)

	return d.repo.Update(ctx, drug)
}

func (d *eatablesService) Delete(ctx context.Context, eatableID string) error {
	return d.repo.Delete(ctx, eatableID)
}

func (d *eatablesService) GetDrugs(ctx context.Context, page, limit uint64, animalID string) (*entity.ListDrugEatables, error) {
	return d.repo.GetDrugs(ctx, page, limit, animalID)
}

func (d *eatablesService) GetFoods(ctx context.Context, page, limit uint64, animalID string) (*entity.ListFoodEatables, error) {
	return d.repo.GetFoods(ctx, page, limit, animalID)
}
