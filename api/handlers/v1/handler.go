package v1

import (
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	"musobaqa/farm-competition/internal/pkg/config"
	tokens "musobaqa/farm-competition/internal/pkg/token"
	"musobaqa/farm-competition/internal/usecase/animals"
	"musobaqa/farm-competition/internal/usecase/category"
	"musobaqa/farm-competition/internal/usecase/drugs"
	"musobaqa/farm-competition/internal/usecase/foods"
	"musobaqa/farm-competition/internal/usecase/products"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Category       category.Category
	Product        products.Product
	Animals        animals.Animal
	Food           foods.Food
	Drug           drugs.Drug
	Enforcer       *casbin.Enforcer
}

type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Category       category.Category
	Product        products.Product
	Animals        animals.Animal
	Food           foods.Food
	Drug           drugs.Drug
	Enforcer       *casbin.Enforcer
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		JwtHandler:     c.JwtHandler,
		Category:       c.Category,
		Product:        c.Product,
		Animals:        c.Animals,
		Food:           c.Food,
		Drug:           c.Drug,
		Enforcer:       c.Enforcer,
	}
}
