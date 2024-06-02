package postgresql

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

type foodRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewFood(db *postgres.PostgresDB) repo.Food {
	return &foodRepo{
		tableName: "foods",
		db:        db,
	}
}

func (a *foodRepo) Create(ctx context.Context, food *entity.Food) (*entity.Food, error) {
	return nil, nil
}

func (a *foodRepo) Update(ctx context.Context, food *entity.Food) (*entity.Food, error) {
	return nil, nil
}

func (a *foodRepo) Delete(ctx context.Context, foodID string) error {
	return nil
}

func (a *foodRepo) Get(ctx context.Context, foodID string) (*entity.Food, error) {
	return nil, nil
}

func (a *foodRepo) List(ctx context.Context, page, limit uint64) (*entity.ListFoods, error) {
	return nil, nil
}
