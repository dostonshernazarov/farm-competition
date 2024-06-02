package postgresql

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

type drugRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewDrug(db *postgres.PostgresDB) repo.Drug {
	return &drugRepo{
		tableName: "drugs",
		db:        db,
	}
}

func (a *drugRepo) Create(ctx context.Context, category *entity.Drug) (*entity.Drug, error) {
	return nil, nil
}

func (a *drugRepo) Update(ctx context.Context, category *entity.Drug) (*entity.Drug, error) {
	return nil, nil
}

func (a *drugRepo) Delete(ctx context.Context, categoryID string) error {
	return nil
}

func (a *drugRepo) Get(ctx context.Context, categoryID string) (*entity.Drug, error) {
	return nil, nil
}

func (a *drugRepo) List(ctx context.Context, page, limit uint64) (*entity.ListDrugs, error) {
	return nil, nil
}
