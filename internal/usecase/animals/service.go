package animals

import (
	"context"
	"github.com/google/uuid"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type animalService struct {
	ctxTimeout time.Duration
	repo       repo.Animal
}

func NewAnimalService(timeout time.Duration, repository repo.Animal) Animal {
	return &animalService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (a *animalService) beforeCreate(animal *entity.Animal) {
	animal.ID = uuid.New().String()
	animal.CreatedAt = time.Now().UTC()
	animal.UpdatedAt = time.Now().UTC()
}

func (a *animalService) beforeUpdate(animal *entity.Animal) {
	animal.UpdatedAt = time.Now().UTC()
}

func (a *animalService) Create(ctx context.Context, animal *entity.Animal) (*entity.Animal, error) {
	a.beforeCreate(animal)

	return a.repo.Create(ctx, animal)
}

func (a *animalService) Update(ctx context.Context, animal *entity.Animal) (*entity.Animal, error) {
	a.beforeUpdate(animal)

	return a.repo.Update(ctx, animal)
}

func (a *animalService) Delete(ctx context.Context, animalID string) error {
	return a.repo.Delete(ctx, animalID)
}

func (a *animalService) Get(ctx context.Context, animalID string) (*entity.Animal, error) {
	return a.repo.Get(ctx, animalID)
}

func (a *animalService) List(ctx context.Context, page, limit uint64) (*entity.ListAnimal, error) {
	return a.repo.List(ctx, page, limit)
}
