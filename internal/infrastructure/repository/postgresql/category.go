package postgresql

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

type categoryRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewCategory(db *postgres.PostgresDB) repo.Category {
	return &categoryRepo{
		tableName: "category",
		db:        db,
	}
}

func (a *categoryRepo) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	return nil, nil
}

func (a *categoryRepo) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	return nil, nil
}

func (a *categoryRepo) Delete(ctx context.Context, categoryID string) error {
	return nil
}

func (a *categoryRepo) Get(ctx context.Context, categoryID string) (*entity.Category, error) {
	return nil, nil
}

func (a *categoryRepo) List(ctx context.Context, page, limit uint64) (*entity.ListCategory, error) {
	return nil, nil
}
