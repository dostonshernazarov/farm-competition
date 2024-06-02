package drugs

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type drugService struct {
	ctxTimeout time.Duration
	repo       repo.Drug
}

func NewDrugService(timeout time.Duration, repository repo.Drug) Drug {
	return &drugService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (d *drugService) beforeCreate(drug *entity.Drug) {
	drug.CreatedAt = time.Now().UTC()
	drug.UpdatedAt = time.Now().UTC()
}

func (d *drugService) beforeUpdate(drug *entity.Drug) {
	drug.UpdatedAt = time.Now().UTC()
}

func (d *drugService) Create(ctx context.Context, drug *entity.Drug) (*entity.Drug, error) {
	d.beforeCreate(drug)

	return d.repo.Create(ctx, drug)
}

func (d *drugService) Update(ctx context.Context, drug *entity.Drug) (*entity.Drug, error) {
	d.beforeUpdate(drug)

	return d.repo.Update(ctx, drug)
}

func (d *drugService) Delete(ctx context.Context, drugID string) error {
	return d.repo.Delete(ctx, drugID)
}

func (d *drugService) Get(ctx context.Context, drugID string) (*entity.Drug, error) {
	return d.repo.Get(ctx, drugID)
}

func (d *drugService) List(ctx context.Context, page, limit uint64) (*entity.ListDrugs, error) {
	return d.repo.List(ctx, page, limit)
}
