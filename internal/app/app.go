package app

import (
	"context"
	"fmt"
	animalproduct "musobaqa/farm-competition/internal/usecase/animal-product"
	"net/http"
	"time"

	"go.uber.org/zap"

	"musobaqa/farm-competition/api"
	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql"
	"musobaqa/farm-competition/internal/pkg/config"
	"musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"musobaqa/farm-competition/internal/pkg/postgres"
	"musobaqa/farm-competition/internal/pkg/redis"

	"musobaqa/farm-competition/internal/usecase/animals"
	"musobaqa/farm-competition/internal/usecase/delivery"
	"musobaqa/farm-competition/internal/usecase/drugs"
	"musobaqa/farm-competition/internal/usecase/foods"
	"musobaqa/farm-competition/internal/usecase/products"
)

type App struct {
	Config        *config.Config
	Logger        *zap.Logger
	DB            *postgres.PostgresDB
	RedisDB       *redis.RedisDB
	server        *http.Server
	ShutdownOTLP  func() error
	Product       products.Product
	Animals       animals.Animal
	Food          foods.Food
	Drug          drugs.Drug
	Delivery      delivery.Delivery
	AnimalProduct animalproduct.AnimalProduct
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// postgres init
	db, err := postgres.New(&cfg)
	if err != nil {
		return nil, err
	}

	// redis init
	redisdb, err := redis.New(&cfg)
	if err != nil {
		return nil, err
	}

	// otlp collector init
	shutdownOTLP, err := otlp.InitOTLPProvider(&cfg)
	if err != nil {
		return nil, err
	}

	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err = time.ParseDuration(cfg.Context.Timeout)
	if err != nil {
		return nil, err
	}

	// product
	productRepo := postgresql.NewProduct(db)
	appProductUseCase := products.NewFoodService(contextTimeout, productRepo)

	// animals
	animalRepo := postgresql.NewAnimal(db)
	appAnimalUseCase := animals.NewAnimalService(contextTimeout, animalRepo)

	// drugs
	drugRepo := postgresql.NewDrug(db)
	appDrugUseCase := drugs.NewDrugService(contextTimeout, drugRepo)

	// food
	foodRepo := postgresql.NewFood(db)
	appFoodUseCase := foods.NewFoodService(contextTimeout, foodRepo)

	// delivery
	deliveryRepo := postgresql.NewDelivery(db)
	appDeliveryUseCase := delivery.NewDeliveryService(contextTimeout, deliveryRepo)

	// animal-product
	animalProductRepo := postgresql.NewAnimalProduct(db)
	appAnimalProductUseCase := animalproduct.NewAnimalProductService(contextTimeout, animalProductRepo)

	return &App{
		Config:        &cfg,
		Logger:        logger,
		DB:            db,
		RedisDB:       redisdb,
		ShutdownOTLP:  shutdownOTLP,
		Product:       appProductUseCase,
		Animals:       appAnimalUseCase,
		Drug:          appDrugUseCase,
		Food:          appFoodUseCase,
		Delivery:      appDeliveryUseCase,
		AnimalProduct: appAnimalProductUseCase,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Product:        a.Product,
		Animals:        a.Animals,
		Food:           a.Food,
		Drug:           a.Drug,
		Delivery:       a.Delivery,
		AnimalProduct:  a.AnimalProduct,
	})
	
	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close database
	a.DB.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
