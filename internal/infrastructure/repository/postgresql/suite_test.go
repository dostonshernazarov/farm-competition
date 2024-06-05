package postgresql_test

import (
	"context"
	"musobaqa/farm-competition/internal/pkg/config"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

func TestAnimalCRUD(t *testing.T) {
	cfg, err := config.NewConfig()
	require.NoError(t, err)
	db, err := postgres.New(cfg)
	require.NoError(t, err)
	defer db.Close()

	repo := postgresql.NewAnimal(db)
	ctx := context.Background()
	animalID := uuid.New().String()

	// create
	animal := &entity.Animal{
		ID:           animalID,
		Name:         "Test Animal",
		CategoryName: "Test Category",
		Gender:       "male",
		BirthDay:     "2023-01-01",
		Genus:        "Test Genus",
		Weight:       100,
		Description:  "Test Description",
		IsHealth:     "true",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	createdAnimal, err := repo.Create(ctx, animal)
	require.NoError(t, err)
	assert.Equal(t, animal.Name, createdAnimal.Name)
	assert.Equal(t, animal.CategoryName, createdAnimal.CategoryName)
	assert.Equal(t, animal.Gender, createdAnimal.Gender)
	assert.Equal(t, animal.BirthDay, createdAnimal.BirthDay)
	assert.Equal(t, animal.Description, createdAnimal.Description)
	assert.Equal(t, animal.IsHealth, createdAnimal.IsHealth)
	assert.Equal(t, animal.Weight, createdAnimal.Weight)

	// update
	updateAnimal := &entity.Animal{
		ID:           animalID,
		Name:         "New Test Animal",
		CategoryName: "Test Category",
		Gender:       "male",
		BirthDay:     "2023-01-01",
		Genus:        "Test Genus",
		Weight:       100,
		Description:  "Test Description",
		IsHealth:     "true",
		UpdatedAt:    time.Now().UTC(),
	}
	updatedAnimal, err := repo.Update(ctx, updateAnimal)
	require.NoError(t, err)
	assert.Equal(t, updateAnimal.Name, updatedAnimal.Name)
	assert.Equal(t, updateAnimal.CategoryName, updatedAnimal.CategoryName)
	assert.Equal(t, updateAnimal.Gender, updatedAnimal.Gender)
	assert.Equal(t, updateAnimal.BirthDay, updatedAnimal.BirthDay)
	assert.Equal(t, updateAnimal.Description, updatedAnimal.Description)
	assert.Equal(t, updateAnimal.IsHealth, updatedAnimal.IsHealth)
	assert.Equal(t, updateAnimal.Weight, updatedAnimal.Weight)

	// get
	getAnimal, err := repo.Get(ctx, animalID)
	require.NoError(t, err)
	assert.Equal(t, updateAnimal.Name, getAnimal.Name)
	assert.Equal(t, updateAnimal.CategoryName, getAnimal.CategoryName)
	assert.Equal(t, updateAnimal.Gender, getAnimal.Gender)
	assert.Equal(t, updateAnimal.BirthDay, getAnimal.BirthDay)
	assert.Equal(t, updateAnimal.Description, getAnimal.Description)
	assert.Equal(t, updateAnimal.IsHealth, getAnimal.IsHealth)
	assert.Equal(t, updateAnimal.Weight, getAnimal.Weight)

	// delete
	err = repo.Delete(ctx, animalID)
	assert.NoError(t, err)
}

func TestProductCRUD(t *testing.T) {
	cfg, err := config.NewConfig()
	require.NoError(t, err)
	db, err := postgres.New(cfg)
	require.NoError(t, err)
	defer db.Close()

	repo := postgresql.NewProduct(db)
	ctx := context.Background()
	productID := uuid.New().String()

	// create
	product := &entity.Product{
		ID:            productID,
		Name:          "Test Product",
		Union:         "Test Union",
		TotalCapacity: 0,
		Description:   "Test Description",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	// update
	updatedProduct, err := repo.Create(ctx, product)
	require.NoError(t, err)
	assert.Equal(t, product.Name, updatedProduct.Name)
	assert.Equal(t, product.Union, updatedProduct.Union)
	assert.Equal(t, product.ID, updatedProduct.ID)
	assert.Equal(t, product.TotalCapacity, updatedProduct.TotalCapacity)

	// get
	getProduct, err := repo.Get(ctx, map[string]string{
		"id": productID,
	})
	assert.Equal(t, getProduct.Description, updatedProduct.Description)
	assert.Equal(t, getProduct.ID, updatedProduct.ID)
	assert.Equal(t, getProduct.TotalCapacity, updatedProduct.TotalCapacity)
	assert.Equal(t, getProduct.TotalCapacity, updatedProduct.TotalCapacity)

	// delete
	err = repo.Delete(ctx, productID)
	assert.NoError(t, err)
}
