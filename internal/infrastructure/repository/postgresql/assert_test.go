package postgresql_test

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/cast"
	"musobaqa/farm-competition/internal/pkg/config"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"musobaqa/farm-competition/internal/entity"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql"
	"musobaqa/farm-competition/internal/pkg/postgres"
)

func TestAnimal(t *testing.T) {
	// Setup Config
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repo := postgresql.NewAnimal(db)
	ctx := context.Background()
	defaultAnimalID := uuid.New().String()

	// Create
	createdAnimalModel := &entity.Animal{
		ID:           defaultAnimalID,
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

	createdAnimal, err := repo.Create(ctx, createdAnimalModel)
	assert.NoError(t, err)
	assert.Equal(t, createdAnimalModel.ID, createdAnimal.ID)
	assert.Equal(t, createdAnimalModel.Name, createdAnimal.Name)
	assert.Equal(t, createdAnimalModel.CategoryName, createdAnimal.CategoryName)
	assert.Equal(t, createdAnimalModel.Gender, createdAnimal.Gender)
	assert.Equal(t, createdAnimal.BirthDay, createdAnimal.BirthDay)
	assert.Equal(t, createdAnimalModel.Description, createdAnimal.Description)
	assert.Equal(t, createdAnimalModel.IsHealth, createdAnimal.IsHealth)
	assert.Equal(t, createdAnimalModel.Weight, createdAnimal.Weight)

	// Update
	updatedAnimalModel := &entity.Animal{
		ID:           defaultAnimalID,
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
	updatedAnimal, err := repo.Update(ctx, updatedAnimalModel)
	assert.NoError(t, err)
	assert.Equal(t, updatedAnimalModel.ID, updatedAnimal.ID)
	assert.Equal(t, updatedAnimalModel.Name, updatedAnimal.Name)
	assert.Equal(t, updatedAnimalModel.CategoryName, updatedAnimal.CategoryName)
	assert.Equal(t, updatedAnimalModel.Gender, updatedAnimal.Gender)
	assert.Equal(t, updatedAnimalModel.BirthDay+"T00:00:00Z", updatedAnimal.BirthDay)
	assert.Equal(t, updatedAnimalModel.Genus, updatedAnimal.Genus)
	assert.Equal(t, updatedAnimalModel.Weight, updatedAnimal.Weight)
	assert.Equal(t, updatedAnimalModel.Description, updatedAnimal.Description)
	assert.Equal(t, updatedAnimalModel.IsHealth, updatedAnimal.IsHealth)
	assert.NotEqual(t, updatedAnimal.Name, createdAnimal.Name)

	// Get
	getAnimal, err := repo.Get(ctx, defaultAnimalID)
	assert.NoError(t, err)
	assert.Equal(t, getAnimal.Name, updatedAnimal.Name)
	assert.Equal(t, getAnimal.CategoryName, updatedAnimal.CategoryName)
	assert.Equal(t, getAnimal.Gender, updatedAnimal.Gender)
	assert.Equal(t, getAnimal.BirthDay, updatedAnimal.BirthDay)
	assert.Equal(t, getAnimal.Genus, updatedAnimal.Genus)
	assert.Equal(t, getAnimal.Weight, updatedAnimal.Weight)
	assert.Equal(t, getAnimal.Description, updatedAnimal.Description)
	assert.Equal(t, getAnimal.IsHealth, updatedAnimal.IsHealth)

	// Delete
	err = repo.Delete(ctx, defaultAnimalID)
	assert.NoError(t, err)
	notAnimal, err := repo.Get(ctx, defaultAnimalID)
	assert.Error(t, err)
	assert.Nil(t, notAnimal)
	assert.Equal(t, err, pgx.ErrNoRows)
}

func TestProduct(t *testing.T) {
	// Setup Config
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repo := postgresql.NewProduct(db)
	ctx := context.Background()
	defaultProductID := uuid.New().String()

	// Create
	createdProductModel := entity.Product{
		ID:            defaultProductID,
		Name:          "Test Product",
		Union:         "Test Union",
		Description:   "Test Description",
		TotalCapacity: 10,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	createdProduct, err := repo.Create(ctx, &createdProductModel)
	assert.NoError(t, err)
	assert.Equal(t, createdProductModel.ID, createdProduct.ID)
	assert.Equal(t, createdProductModel.Name, createdProduct.Name)
	assert.Equal(t, createdProductModel.Union, createdProduct.Union)
	assert.Equal(t, createdProductModel.Description, createdProduct.Description)
	assert.Equal(t, createdProductModel.TotalCapacity, createdProduct.TotalCapacity)

	updatedProductModel := entity.Product{
		ID:            defaultProductID,
		Name:          "New Test Product",
		Union:         "Test Union",
		Description:   "Test Description",
		TotalCapacity: 10,
		UpdatedAt:     time.Now().UTC(),
	}
	updatedProduct, err := repo.Update(ctx, &updatedProductModel)
	assert.NoError(t, err)
	assert.Equal(t, updatedProductModel.ID, updatedProduct.ID)
	assert.Equal(t, updatedProductModel.Name, updatedProduct.Name)
	assert.Equal(t, updatedProductModel.Union, updatedProduct.Union)
	assert.Equal(t, updatedProductModel.Description, updatedProduct.Description)
	assert.Equal(t, updatedProductModel.TotalCapacity, updatedProduct.TotalCapacity)
	assert.NotEqual(t, updatedProduct.Name, createdProduct.Name)

	// Get
	getProduct, err := repo.Get(ctx, map[string]string{
		"id": defaultProductID,
	})
	assert.NoError(t, err)
	assert.Equal(t, getProduct.ID, updatedProduct.ID)
	assert.Equal(t, getProduct.Name, updatedProduct.Name)
	assert.Equal(t, getProduct.Union, updatedProduct.Union)
	assert.Equal(t, getProduct.Description, updatedProduct.Description)
	assert.Equal(t, getProduct.TotalCapacity, updatedProduct.TotalCapacity)

	// delete
	err = repo.Delete(ctx, defaultProductID)
	assert.NoError(t, err)
	notProduct, err := repo.Get(ctx, map[string]string{
		"id": defaultProductID,
	})
	assert.Error(t, err)
	assert.Nil(t, notProduct)
	assert.Equal(t, err, pgx.ErrNoRows)
}

func TestAnimalProduct(t *testing.T) {
	// Setup Config
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repoAnimal := postgresql.NewAnimal(db)
	repoProduct := postgresql.NewProduct(db)
	repoAnimalProduct := postgresql.NewAnimalProduct(db)
	ctx := context.Background()

	// Default Values
	defaultProductID := uuid.New().String()
	defaultAnimalID := uuid.New().String()
	defaultAnimalProductID := uuid.New().String()

	// Create Animal
	createdAnimalModel := entity.Animal{
		ID:           defaultAnimalID,
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

	createdAnimal, err := repoAnimal.Create(ctx, &createdAnimalModel)
	assert.NoError(t, err)

	// Create Product
	createdProductModel := entity.Product{
		ID:            defaultProductID,
		Name:          "Test Product",
		Union:         "Test Union",
		Description:   "Test Description",
		TotalCapacity: 10,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	createdProduct, err := repoProduct.Create(ctx, &createdProductModel)
	assert.NoError(t, err)

	// Create Animal-Product
	createdAnimalProductModel := entity.AnimalProductReq{
		ID:        defaultAnimalProductID,
		AnimalID:  defaultAnimalID,
		ProductID: defaultProductID,
		Capacity:  10,
		GetTime:   "2024-06-05 14:00:00",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// "2024-06-05 14:00:00" -> "2024-06-05T14:00:00Z" // Database CURRENT_TIMESTAMP format
	getTimeArray := strings.Split(createdAnimalProductModel.GetTime, " ")
	getTime := getTimeArray[0] + "T" + getTimeArray[1] + "Z"

	createdAnimalProduct, err := repoAnimalProduct.Create(ctx, &createdAnimalProductModel)
	assert.NoError(t, err)
	assert.Equal(t, createdAnimalProductModel.ID, createdAnimalProduct.ID)
	assert.Equal(t, createdAnimalProductModel.AnimalID, createdAnimalProduct.Animal.ID)
	assert.Equal(t, createdAnimalProductModel.ProductID, createdAnimalProduct.Product.ID)
	assert.Equal(t, getTime, createdAnimalProduct.GetTime)
	assert.Equal(t, createdAnimalProductModel.Capacity, createdAnimalProduct.Capacity)
	assert.Equal(t, createdAnimalProduct.Animal.ID, createdAnimal.ID)
	assert.Equal(t, createdAnimalProduct.Animal.ID, defaultAnimalID)
	assert.Equal(t, createdAnimalProduct.Animal.Name, createdAnimal.Name)
	assert.Equal(t, createdAnimalProduct.Animal.CategoryName, createdAnimal.CategoryName)
	assert.Equal(t, createdAnimalProduct.Animal.Description, createdAnimal.Description)
	assert.Equal(t, createdAnimalProduct.Animal.Gender, createdAnimal.Gender)
	assert.Equal(t, createdAnimalProduct.Animal.BirthDay, createdAnimal.BirthDay)
	assert.Equal(t, createdAnimalProduct.Animal.Genus, createdAnimal.Genus)
	assert.Equal(t, createdAnimalProduct.Animal.Weight, createdAnimal.Weight)
	assert.Equal(t, createdAnimalProduct.Animal.Description, createdAnimal.Description)
	assert.Equal(t, createdAnimalProduct.Animal.IsHealth, createdAnimal.IsHealth)
	assert.Equal(t, createdAnimalProduct.Product.ID, createdProduct.ID)
	assert.Equal(t, createdAnimalProduct.Product.ID, defaultProductID)
	assert.Equal(t, createdAnimalProduct.Product.Name, createdProduct.Name)
	assert.Equal(t, createdAnimalProduct.Product.Union, createdProduct.Union)
	assert.Equal(t, createdAnimalProduct.Product.Description, createdProduct.Description)
	assert.Equal(t, createdAnimalProduct.Product.TotalCapacity, createdProduct.TotalCapacity)

	// Update Animal-Product
	updatedAnimalProductModel := entity.AnimalProductReq{
		ID:        defaultAnimalProductID,
		AnimalID:  defaultAnimalID,
		ProductID: defaultProductID,
		Capacity:  11,
		GetTime:   "2024-06-05 14:00:00",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// "2024-06-05 14:00:00" -> "2024-06-05T14:00:00Z" // Database CURRENT_TIMESTAMP format
	getTimeArray = strings.Split(updatedAnimalProductModel.GetTime, " ")
	getTime = getTimeArray[0] + "T" + getTimeArray[1] + "Z"

	updatedAnimalProduct, err := repoAnimalProduct.Update(ctx, &updatedAnimalProductModel)
	assert.NoError(t, err)
	assert.Equal(t, updatedAnimalProductModel.ID, updatedAnimalProduct.ID)
	assert.Equal(t, updatedAnimalProductModel.AnimalID, updatedAnimalProduct.Animal.ID)
	assert.Equal(t, updatedAnimalProductModel.ProductID, updatedAnimalProduct.Product.ID)
	assert.Equal(t, getTime, updatedAnimalProduct.GetTime)
	assert.NotEqual(t, updatedAnimalProduct.Capacity, createdAnimalProduct.Capacity)
	assert.Equal(t, updatedAnimalProduct.Capacity, int64(11))
	assert.Equal(t, updatedAnimalProduct.Animal.ID, createdAnimal.ID)
	assert.Equal(t, updatedAnimalProduct.Animal.ID, defaultAnimalID)
	assert.Equal(t, updatedAnimalProduct.Animal.Name, createdAnimal.Name)
	assert.Equal(t, updatedAnimalProduct.Animal.CategoryName, createdAnimal.CategoryName)
	assert.Equal(t, updatedAnimalProduct.Animal.Description, createdAnimal.Description)
	assert.Equal(t, updatedAnimalProduct.Animal.Gender, createdAnimal.Gender)
	assert.Equal(t, updatedAnimalProduct.Animal.BirthDay, createdAnimal.BirthDay)
	assert.Equal(t, updatedAnimalProduct.Animal.Genus, createdAnimal.Genus)
	assert.Equal(t, updatedAnimalProduct.Animal.Weight, createdAnimal.Weight)
	assert.Equal(t, updatedAnimalProduct.Animal.Description, createdAnimal.Description)
	assert.Equal(t, updatedAnimalProduct.Animal.IsHealth, createdAnimal.IsHealth)
	assert.Equal(t, updatedAnimalProduct.Product.ID, createdProduct.ID)
	assert.Equal(t, updatedAnimalProduct.Product.ID, defaultProductID)
	assert.Equal(t, updatedAnimalProduct.Product.Name, createdProduct.Name)
	assert.Equal(t, updatedAnimalProduct.Product.Union, createdProduct.Union)
	assert.Equal(t, updatedAnimalProduct.Product.Description, createdProduct.Description)
	assert.Equal(t, updatedAnimalProduct.Product.TotalCapacity, createdProduct.TotalCapacity)

	// Get Animal-Product
	getAnimalProduct, err := repoAnimalProduct.Get(ctx, defaultAnimalProductID)
	assert.NoError(t, err)
	assert.Equal(t, updatedAnimalProduct.ID, getAnimalProduct.ID)
	assert.Equal(t, updatedAnimalProduct.Animal.ID, getAnimalProduct.Animal.ID)
	assert.Equal(t, updatedAnimalProduct.Product.ID, getAnimalProduct.Product.ID)
	assert.Equal(t, getTime, getAnimalProduct.GetTime)
	assert.Equal(t, updatedAnimalProduct.Capacity, getAnimalProduct.Capacity)
	assert.Equal(t, getAnimalProduct.Product.ID, defaultProductID)
	assert.Equal(t, getAnimalProduct.Animal.ID, defaultAnimalID)
	assert.Equal(t, getAnimalProduct.Animal.Name, updatedAnimalProduct.Animal.Name)
	assert.Equal(t, getAnimalProduct.Animal.CategoryName, updatedAnimalProduct.Animal.CategoryName)
	assert.Equal(t, getAnimalProduct.Animal.Description, updatedAnimalProduct.Animal.Description)
	assert.Equal(t, getAnimalProduct.Animal.Gender, updatedAnimalProduct.Animal.Gender)
	assert.Equal(t, getAnimalProduct.Animal.BirthDay, updatedAnimalProduct.Animal.BirthDay)
	assert.Equal(t, getAnimalProduct.Animal.Genus, updatedAnimalProduct.Animal.Genus)
	assert.Equal(t, getAnimalProduct.Animal.Weight, updatedAnimalProduct.Animal.Weight)
	assert.Equal(t, getAnimalProduct.Animal.Description, updatedAnimalProduct.Animal.Description)
	assert.Equal(t, getAnimalProduct.Animal.IsHealth, updatedAnimalProduct.Animal.IsHealth)
	assert.Equal(t, getAnimalProduct.Product.Name, updatedAnimalProduct.Product.Name)
	assert.Equal(t, getAnimalProduct.Product.Union, updatedAnimalProduct.Product.Union)
	assert.Equal(t, getAnimalProduct.Product.Description, updatedAnimalProduct.Product.Description)
	assert.Equal(t, getAnimalProduct.Product.TotalCapacity, updatedAnimalProduct.Product.TotalCapacity)

	// Delete Animal-Product
	err = repoAnimalProduct.Delete(ctx, defaultAnimalProductID)
	assert.NoError(t, err)
	notAnimalProduct, err := repoAnimalProduct.Get(ctx, defaultAnimalProductID)
	assert.Error(t, err)
	assert.Nil(t, notAnimalProduct)
	assert.Equal(t, err, pgx.ErrNoRows)
}

func TestFood(t *testing.T) {
	// Setup Config
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repo := postgresql.NewFood(db)
	ctx := context.Background()
	defaultFoodID := uuid.New().String()

	// Create
	createdFoodModel := entity.Food{
		ID:          defaultFoodID,
		Name:        "Test Food",
		Capacity:    10,
		Union:       "kilogram",
		Description: "For all animals",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	createdFood, err := repo.Create(ctx, &createdFoodModel)
	assert.NoError(t, err)
	assert.Equal(t, createdFood.ID, defaultFoodID)
	assert.Equal(t, createdFood.ID, createdFoodModel.ID)
	assert.Equal(t, createdFood.Name, createdFoodModel.Name)
	assert.Equal(t, createdFood.Capacity, createdFoodModel.Capacity)
	assert.Equal(t, createdFood.Union, createdFoodModel.Union)
	assert.Equal(t, createdFood.Description, createdFoodModel.Description)

	// Update
	updatedFoodModel := entity.Food{
		ID:          defaultFoodID,
		Name:        "Updated Food",
		Capacity:    10,
		Union:       "kilogram",
		Description: "For all animals",
		UpdatedAt:   time.Now().UTC(),
	}

	updatedFood, err := repo.Update(ctx, &updatedFoodModel)
	assert.NoError(t, err)
	assert.Equal(t, updatedFood.ID, defaultFoodID)
	assert.Equal(t, updatedFood.ID, updatedFoodModel.ID)
	assert.Equal(t, updatedFood.Name, updatedFoodModel.Name)
	assert.Equal(t, updatedFood.Capacity, updatedFoodModel.Capacity)
	assert.Equal(t, updatedFood.Union, updatedFoodModel.Union)
	assert.Equal(t, updatedFood.Description, updatedFoodModel.Description)

	// Get
	getFood, err := repo.Get(ctx, map[string]string{
		"id": defaultFoodID,
	})
	assert.NoError(t, err)
	assert.Equal(t, getFood.ID, defaultFoodID)
	assert.Equal(t, getFood.Name, updatedFood.Name)
	assert.Equal(t, getFood.Capacity, updatedFood.Capacity)
	assert.Equal(t, getFood.Union, updatedFood.Union)
	assert.Equal(t, getFood.Description, updatedFood.Description)

	// Delete
	err = repo.Delete(ctx, defaultFoodID)
	assert.NoError(t, err)
	notFood, err := repo.Get(ctx, map[string]string{
		"id": defaultFoodID,
	})
	assert.Error(t, err)
	assert.Nil(t, notFood)
	assert.Equal(t, err, pgx.ErrNoRows)
}

func TestDrug(t *testing.T) {
	// Setup Config
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repo := postgresql.NewDrug(db)
	ctx := context.Background()
	defaultDrugID := uuid.New().String()

	// Create
	createdDrugModel := entity.Drug{
		ID:          defaultDrugID,
		Name:        "Test Food",
		Capacity:    10,
		Union:       "kilogram",
		Description: "For all animals",
		Status:      "test Status",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	createdDrug, err := repo.Create(ctx, &createdDrugModel)
	assert.NoError(t, err)
	assert.Equal(t, createdDrug.ID, defaultDrugID)
	assert.Equal(t, createdDrug.ID, createdDrugModel.ID)
	assert.Equal(t, createdDrug.Name, createdDrugModel.Name)
	assert.Equal(t, createdDrug.Capacity, createdDrugModel.Capacity)
	assert.Equal(t, createdDrug.Union, createdDrugModel.Union)
	assert.Equal(t, createdDrug.Description, createdDrugModel.Description)

	// Update
	updatedDrugModel := entity.Drug{
		ID:          defaultDrugID,
		Name:        "Updated Food",
		Capacity:    10,
		Union:       "kilogram",
		Description: "For all animals",
		Status:      "Test Status",
		UpdatedAt:   time.Now().UTC(),
	}

	updatedDrug, err := repo.Update(ctx, &updatedDrugModel)
	assert.NoError(t, err)
	assert.Equal(t, updatedDrug.ID, defaultDrugID)
	assert.Equal(t, updatedDrug.ID, updatedDrugModel.ID)
	assert.Equal(t, updatedDrug.Name, updatedDrugModel.Name)
	assert.Equal(t, updatedDrug.Capacity, updatedDrugModel.Capacity)
	assert.Equal(t, updatedDrug.Union, updatedDrugModel.Union)
	assert.Equal(t, updatedDrug.Description, updatedDrugModel.Description)

	// Get
	getDrug, err := repo.Get(ctx, map[string]string{
		"id": defaultDrugID,
	})
	assert.NoError(t, err)
	assert.Equal(t, getDrug.ID, defaultDrugID)
	assert.Equal(t, getDrug.Name, updatedDrug.Name)
	assert.Equal(t, getDrug.Capacity, updatedDrug.Capacity)
	assert.Equal(t, getDrug.Union, updatedDrug.Union)
	assert.Equal(t, getDrug.Description, updatedDrug.Description)

	// Delete
	err = repo.Delete(ctx, defaultDrugID)
	assert.NoError(t, err)
	notFood, err := repo.Get(ctx, map[string]string{
		"id": defaultDrugID,
	})
	assert.Error(t, err)
	assert.Nil(t, notFood)
	assert.Equal(t, err, pgx.ErrNoRows)
}

func TestEatable(t *testing.T) {
	// Setup Config
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repoEatable := postgresql.NewEatable(db)
	repoAnimal := postgresql.NewAnimal(db)
	repoFood := postgresql.NewFood(db)

	// Setup Default Values
	ctx := context.Background()
	defaultEatableID := uuid.New().String()
	defaultAnimalID := uuid.New().String()
	defaultFoodID := uuid.New().String()

	// Create Animal
	createdAnimalModel := entity.Animal{
		ID:           defaultAnimalID,
		Name:         "Test Animal",
		CategoryName: "Test Category",
		Gender:       "male",
		BirthDay:     "2024-06-06",
		Genus:        "Test Genus",
		Weight:       10,
		IsHealth:     "true",
		Description:  "Test Description",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	createdAnimal, err := repoAnimal.Create(ctx, &createdAnimalModel)
	assert.NoError(t, err)
	// Create Food
	createdFoodModel := entity.Food{
		ID:          defaultFoodID,
		Name:        "Test Animal",
		Capacity:    10,
		Union:       "kilogram",
		Description: "Test Description",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	createdFood, err := repoFood.Create(ctx, &createdFoodModel)
	assert.NoError(t, err)

	// Make JSONB format
	var daily []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	}

	for i := 0; i < 3; i++ {
		daily = append(daily, struct {
			Capacity int64  `json:"capacity"`
			Time     string `json:"time"`
		}{
			Capacity: int64(i + 5),
			Time:     cast.ToString(i+10) + ":00:00",
		})
	}

	// Create Eatables
	createdEatablesModel := entity.Eatables{
		ID:        defaultEatableID,
		AnimalID:  createdAnimal.ID,
		EatableID: createdFood.ID,
		Category:  "food",
		Daily:     daily,
	}

	createdEatables, err := repoEatable.Create(ctx, &createdEatablesModel)
	assert.NoError(t, err)
	assert.Equal(t, createdEatables.ID, defaultEatableID)
	assert.Equal(t, createdEatables.ID, createdEatablesModel.ID)
	assert.Equal(t, createdEatables.Eatable.ID, defaultFoodID)
	assert.Equal(t, createdEatables.Eatable.ID, createdEatablesModel.EatableID)
	assert.Equal(t, createdEatables.AnimalID, defaultAnimalID)
	assert.Equal(t, createdEatables.AnimalID, createdEatablesModel.AnimalID)
	assert.Equal(t, createdEatables.Category, createdEatablesModel.Category)
	assert.Equal(t, createdEatables.Eatable.Union, createdFood.Union)
	assert.Equal(t, createdEatables.Eatable.Capacity, createdFood.Capacity)
	assert.Equal(t, createdEatables.Eatable.Name, createdFood.Name)
	assert.Equal(t, createdEatables.Eatable.Description, createdFood.Description)
	assert.Equal(t, createdEatables.Daily[0].Time, daily[0].Time)
	assert.Equal(t, createdEatables.Daily[0].Capacity, daily[0].Capacity)
	assert.Equal(t, createdEatables.Daily[1].Time, daily[1].Time)
	assert.Equal(t, createdEatables.Daily[1].Capacity, daily[1].Capacity)
	assert.Equal(t, createdEatables.Daily[2].Time, daily[2].Time)
	assert.Equal(t, createdEatables.Daily[2].Capacity, daily[2].Capacity)

	// Update Eatable
	updatedEatablesModel := entity.Eatables{
		ID:        defaultEatableID,
		AnimalID:  defaultAnimalID,
		EatableID: defaultFoodID,
		Category:  "food",
		Daily:     daily[:len(daily)-1],
	}

	updatedEatables, err := repoEatable.Update(ctx, &updatedEatablesModel)
	assert.NoError(t, err)
	assert.Equal(t, updatedEatables.ID, defaultEatableID)
	assert.Equal(t, updatedEatables.ID, updatedEatablesModel.ID)
	assert.Equal(t, updatedEatables.Eatable.ID, defaultFoodID)
	assert.Equal(t, updatedEatables.Eatable.ID, updatedEatablesModel.EatableID)
	assert.Equal(t, updatedEatables.Category, updatedEatablesModel.Category)
	assert.Equal(t, updatedEatables.Daily[0].Time, updatedEatablesModel.Daily[0].Time)
	assert.Equal(t, updatedEatables.Daily[0].Capacity, updatedEatablesModel.Daily[0].Capacity)
	assert.Equal(t, updatedEatables.Daily[1].Time, updatedEatablesModel.Daily[1].Time)
	assert.Equal(t, updatedEatables.Daily[1].Capacity, updatedEatablesModel.Daily[1].Capacity)

	// Delete Eatables
	err = repoEatable.Delete(ctx, defaultEatableID)
	assert.NoError(t, err)
}

func TestDelivery(t *testing.T) {
	// Setup Config
	cfg, err := config.NewConfig()
	assert.NoError(t, err)

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repo := postgresql.NewDelivery(db)
	ctx := context.Background()
	defaultDeliveryID := uuid.New().String()

	// Create
	createdDeliveryModel := &entity.Delivery{
		ID:        defaultDeliveryID,
		Name:      "Test Product",
		Union:     "Test Union",
		Category:  "Test category",
		Capacity:  10,
		Time:      "2024-06-06 13:00:00",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdDelivery, err := repo.Create(ctx, createdDeliveryModel)
	assert.NoError(t, err)
	assert.Equal(t, createdDelivery.ID, defaultDeliveryID)
	assert.Equal(t, createdDelivery.ID, createdDeliveryModel.ID)
	assert.Equal(t, createdDelivery.Name, createdDeliveryModel.Name)
	assert.Equal(t, createdDelivery.Union, createdDeliveryModel.Union)
	assert.Equal(t, createdDelivery.Capacity, createdDeliveryModel.Capacity)
	assert.Equal(t, createdDelivery.Time, createdDeliveryModel.Time)

	// Update
	updatedDeliveryModel := &entity.Delivery{
		ID:        defaultDeliveryID,
		Name:      "New Test Product",
		Union:     "Test Union",
		Category:  "Test category",
		Capacity:  10,
		Time:      "2024-06-06 13:00:00",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	updatedDelivery, err := repo.Update(ctx, updatedDeliveryModel)
	assert.NoError(t, err)
	assert.Equal(t, updatedDelivery.ID, defaultDeliveryID)
	assert.Equal(t, updatedDelivery.ID, updatedDeliveryModel.ID)
	assert.Equal(t, updatedDelivery.Name, updatedDeliveryModel.Name)
	assert.Equal(t, updatedDelivery.Category, createdDeliveryModel.Category)
	assert.Equal(t, updatedDelivery.Capacity, updatedDeliveryModel.Capacity)
	assert.Equal(t, updatedDelivery.Union, updatedDeliveryModel.Union)
	assert.Equal(t, updatedDelivery.Time, updatedDeliveryModel.Time)

	// Get
	getTimeArray := strings.Split(updatedDelivery.Time, " ")
	getTime := getTimeArray[0] + "T" + getTimeArray[1] + "Z"

	getDelivery, err := repo.Get(ctx, defaultDeliveryID)
	assert.NoError(t, err)
	assert.Equal(t, getDelivery.ID, defaultDeliveryID)
	assert.Equal(t, getDelivery.ID, updatedDelivery.ID)
	assert.Equal(t, getDelivery.Name, updatedDelivery.Name)
	assert.Equal(t, getDelivery.Category, updatedDelivery.Category)
	assert.Equal(t, getDelivery.Capacity, updatedDelivery.Capacity)
	assert.Equal(t, getDelivery.Union, updatedDelivery.Union)
	assert.Equal(t, getDelivery.Time, getTime)

	// Delete
	err = repo.Delete(ctx, defaultDeliveryID)
	assert.NoError(t, err)
	notDelivery, err := repo.Get(ctx, defaultDeliveryID)
	assert.Error(t, err)
	assert.Nil(t, notDelivery)
	assert.Equal(t, err, pgx.ErrNoRows)
}
