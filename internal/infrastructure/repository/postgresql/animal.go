package postgresql

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

type animalRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewAnimal(db *postgres.PostgresDB) repo.Animal {
	return &animalRepo{
		tableName: "animals",
		db:        db,
	}
}

func (a *animalRepo) Create(ctx context.Context, animal *entity.Animal) (*entity.Animal, error) {
	return nil, nil
}

func (a *animalRepo) Update(ctx context.Context, animal *entity.Animal) (*entity.Animal, error) {
	return nil, nil
}

func (a *animalRepo) Delete(ctx context.Context, animalID string) error {
	return nil
}

func (a *animalRepo) Get(ctx context.Context, animalID string) (*entity.Animal, error) {
	return nil, nil
}

func (a *animalRepo) List(ctx context.Context, page, limit uint64) (*entity.ListAnimal, error) {
	return nil, nil
}
