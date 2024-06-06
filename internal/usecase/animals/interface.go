package animals

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Animal interface {
	Create(ctx context.Context, animal *entity.Animal) (*entity.Animal, error)
	Update(ctx context.Context, animal *entity.Animal) (*entity.Animal, error)
	Delete(ctx context.Context, animalID string) error
	Get(ctx context.Context, animalID string) (*entity.Animal, error)
	List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListAnimal, error)
	HungryAnimals(ctx context.Context, page, limit uint64) (*entity.ListAnimal, error)
}
