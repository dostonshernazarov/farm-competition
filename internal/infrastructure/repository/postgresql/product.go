package postgresql

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

type productRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewProduct(db *postgres.PostgresDB) repo.Product {
	return &productRepo{
		tableName: "products",
		db:        db,
	}
}

func (a *productRepo) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	return nil, nil
}

func (a *productRepo) Update(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	return nil, nil
}

func (a *productRepo) Delete(ctx context.Context, productID string) error {
	return nil
}

func (a *productRepo) Get(ctx context.Context, productID string) (*entity.Product, error) {
	return nil, nil
}

func (a *productRepo) List(ctx context.Context, page, limit uint64) (*entity.ListProducts, error) {
	return nil, nil
}
