package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	"musobaqa/farm-competition/api"

	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"

	"musobaqa/farm-competition/internal/infrastructure/repository/postgresql"
	"musobaqa/farm-competition/internal/pkg/config"
	"musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"

	"musobaqa/farm-competition/internal/pkg/postgres"
	"musobaqa/farm-competition/internal/pkg/redis"
	"musobaqa/farm-competition/internal/usecase/animals"
	animlaService "musobaqa/farm-competition/internal/usecase/animals"
	"musobaqa/farm-competition/internal/usecase/category"
	categoryService "musobaqa/farm-competition/internal/usecase/category"
	"musobaqa/farm-competition/internal/usecase/drugs"
	drugService "musobaqa/farm-competition/internal/usecase/drugs"
	"musobaqa/farm-competition/internal/usecase/foods"
	foodService "musobaqa/farm-competition/internal/usecase/foods"
	"musobaqa/farm-competition/internal/usecase/products"
	productService "musobaqa/farm-competition/internal/usecase/products"
)

type Services struct {
	Category category.Category
	Product  products.Product
	Animals  animals.Animal
	Food     foods.Food
	Drug     drugs.Drug
}

type App struct {
	Config       *config.Config
	Logger       *zap.Logger
	DB           *postgres.PostgresDB
	RedisDB      *redis.RedisDB
	server       *http.Server
	Enforcer     *casbin.Enforcer
	ShutdownOTLP func() error
	services     Services
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

	// enforcer.SetCache(policy.NewCache(&redisdb.Client))

	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err = time.ParseDuration(cfg.Context.Timeout)
	if err != nil {
		return nil, err
	}

	// category
	categoryRepo := postgresql.NewCategory(db)
	appCategoryUseCase := categoryService.NewCategoryService(contextTimeout, categoryRepo)

	// product
	productRepo := postgresql.NewProduct(db)
	appProductUseCase := productService.NewFoodService(contextTimeout, productRepo)

	// animals
	animalRepo := postgresql.NewAnimal(db)
	appAnimalUseCase := animlaService.NewAnimalService(contextTimeout, animalRepo)

	// drugs
	drugRepo := postgresql.NewDrug(db)
	appDrugUseCase := drugService.NewDrugService(contextTimeout, drugRepo)

	// food
	foodRepo := postgresql.NewFood(db)
	appFoodUseCase := foodService.NewFoodService(contextTimeout, foodRepo)

	services := Services{
		Category: appCategoryUseCase,
		Product:  appProductUseCase,
		Animals:  appAnimalUseCase,
		Food:     appFoodUseCase,
		Drug:     appDrugUseCase,
	}

	return &App{
		Config:   &cfg,
		Logger:   logger,
		DB:       db,
		RedisDB:  redisdb,
		Enforcer: enforcer,
		// BrokerProducer: kafkaProducer,
		ShutdownOTLP: shutdownOTLP,
		services:     services,
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
		// Cache:          cache,
		Enforcer: a.Enforcer,
		//AppVersion: a.appVersion,
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
