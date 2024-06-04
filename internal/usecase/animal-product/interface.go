package animalproduct

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type AnimalProduct interface {
	Create(ctx context.Context, animalProduct *entity.AnimalProductReq) (*entity.AnimalProductRes, error)
	Update(ctx context.Context, animalProduct *entity.AnimalProductReq) (*entity.AnimalProductRes, error)
	Delete(ctx context.Context, animalProductID string) error
	Get(ctx context.Context, animalProductID string) (*entity.AnimalProductRes, error)
	List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListAnimalProduct, error)
	ListAnimals(ctx context.Context, page, limit uint64, productID string) (*entity.AnimalsWithProduct, error)
	ListProducts(ctx context.Context, page, limit uint64, animalID string) (*entity.ProductsWithAnimal, error)
}
