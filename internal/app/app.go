package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"

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
	Config       *config.Config
	Logger       *zap.Logger
	DB           *postgres.PostgresDB
	RedisDB      *redis.RedisDB
	server       *http.Server
	Enforcer     *casbin.Enforcer
	ShutdownOTLP func() error
	Product      products.Product
	Animals      animals.Animal
	Food         foods.Food
	Drug         drugs.Drug
	Delivery     delivery.Delivery
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

	// initialization enforcer
	enforcer, err := casbin.NewEnforcer("auth.conf", "auth.csv")
	if err != nil {
		return nil, err
	}

	//enforcer.SetCache(policy.NewCache(&redisdb.Client))

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

	return &App{
		Config:       &cfg,
		Logger:       logger,
		DB:           db,
		RedisDB:      redisdb,
		Enforcer:     enforcer,
		ShutdownOTLP: shutdownOTLP,
		Product:      appProductUseCase,
		Animals:      appAnimalUseCase,
		Drug:         appDrugUseCase,
		Food:         appFoodUseCase,
		Delivery:     appDeliveryUseCase,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	// initialize cache
	// cache := redisrepo.NewCache(a.RedisDB)

	// tokenRepo := postgresql.NewRefreshTokenRepo(a.DB)

	// initialize token service
	// refreshTokenService := refresh_token.NewRefreshTokenService(contextTimeout, tokenRepo)

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Enforcer:       a.Enforcer,
		Product:        a.Product,
		Animals:        a.Animals,
		Food:           a.Food,
		Drug:           a.Drug,
		Delivery:       a.Delivery,
	})
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	roleManager := a.Enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl)

	roleManager.AddMatchingFunc("keyMatch", util.KeyMatch)
	roleManager.AddMatchingFunc("keyMatch3", util.KeyMatch3)

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
