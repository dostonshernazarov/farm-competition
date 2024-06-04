package animalproduct

import (
	"context"
	"github.com/google/uuid"
	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql/repo"
	"time"
)

type animalProductService struct {
	ctxTimeout time.Duration
	repo       repo.AnimalProduct
}

func NewAnimalProductService(timeout time.Duration, repository repo.AnimalProduct) AnimalProduct {
	return &animalProductService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (ap *animalProductService) beforeCreate(animal *entity.AnimalProductReq) {
	animal.ID = uuid.New().String()
	animal.CreatedAt = time.Now().UTC()
	animal.UpdatedAt = time.Now().UTC()
}

func (ap *animalProductService) beforeUpdate(animal *entity.AnimalProductReq) {
	animal.UpdatedAt = time.Now().UTC()
}

func (ap *animalProductService) Create(ctx context.Context, animal *entity.AnimalProductReq) (*entity.AnimalProductRes, error) {
	ap.beforeCreate(animal)

	return ap.repo.Create(ctx, animal)
}

func (ap *animalProductService) Update(ctx context.Context, animalProduct *entity.AnimalProductReq) (*entity.AnimalProductRes, error) {
	ap.beforeUpdate(animalProduct)

	return ap.repo.Update(ctx, animalProduct)
}

func (ap *animalProductService) Delete(ctx context.Context, animalProductID string) error {
	return ap.repo.Delete(ctx, animalProductID)
}

func (ap *animalProductService) Get(ctx context.Context, animalProductID string) (*entity.AnimalProductRes, error) {
	return ap.repo.Get(ctx, animalProductID)
}

func (ap *animalProductService) List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListAnimalProduct, error) {
	return ap.repo.List(ctx, page, limit, params)
}

func (ap *animalProductService) ListAnimals(ctx context.Context, page, limit uint64, productID string) (*entity.AnimalsWithProduct, error) {
	return ap.repo.ListAnimals(ctx, page, limit, productID)
}

func (ap *animalProductService) ListProducts(ctx context.Context, page, limit uint64, productID string) (*entity.ProductsWithAnimal, error) {
	return ap.repo.ListProducts(ctx, page, limit, productID)
}
