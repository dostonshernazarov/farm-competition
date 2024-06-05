package v1

import (
	"time"

	"go.uber.org/zap"

	"musobaqa/farm-competition/internal/pkg/config"
	tokens "musobaqa/farm-competition/internal/pkg/token"
	animalproduct "musobaqa/farm-competition/internal/usecase/animal-product"
	"musobaqa/farm-competition/internal/usecase/animals"
	"musobaqa/farm-competition/internal/usecase/delivery"
	"musobaqa/farm-competition/internal/usecase/drugs"
	"musobaqa/farm-competition/internal/usecase/eatables"
	"musobaqa/farm-competition/internal/usecase/feeding"
	"musobaqa/farm-competition/internal/usecase/foods"
	"musobaqa/farm-competition/internal/usecase/products"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Product        products.Product
	Animals        animals.Animal
	Food           foods.Food
	Drug           drugs.Drug
	Delivery       delivery.Delivery
	AnimalProduct  animalproduct.AnimalProduct
	EatablesInfo eatables.Eatable
	Feeding feeding.Feeding
}

type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Product        products.Product
	Animals        animals.Animal
	Food           foods.Food
	Drug           drugs.Drug
	Delivery       delivery.Delivery
	AnimalProduct  animalproduct.AnimalProduct
	EatablesInfo eatables.Eatable
	Feeding feeding.Feeding
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		JwtHandler:     c.JwtHandler,
		Product:        c.Product,
		Animals:        c.Animals,
		Food:           c.Food,
		Drug:           c.Drug,
		Delivery:       c.Delivery,
		AnimalProduct:  c.AnimalProduct,
		EatablesInfo: c.EatablesInfo,
		Feeding: c.Feeding,
	}
}
