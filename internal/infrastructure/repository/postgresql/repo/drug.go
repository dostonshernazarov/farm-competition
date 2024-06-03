package repo

import (
	"context"
	"musobaqa/farm-competition/internal/entity"
)

type Drug interface {
	Create(ctx context.Context, drug *entity.Drug) (*entity.Drug, error)
	Update(ctx context.Context, drug *entity.Drug) (*entity.Drug, error)
	Delete(ctx context.Context, drugID string) error
	Get(ctx context.Context, drugID string) (*entity.Drug, error)
	List(ctx context.Context, page, limit uint64, params map[string]any) (*entity.ListDrugs, error)
	UniqueDrugName(ctx context.Context, drugName string) (int, error)
}
