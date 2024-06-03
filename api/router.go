package api

import (
	"musobaqa/farm-competition/internal/usecase/animals"
	"musobaqa/farm-competition/internal/usecase/drugs"
	"musobaqa/farm-competition/internal/usecase/foods"
	"musobaqa/farm-competition/internal/usecase/products"
	"time"

	_ "musobaqa/farm-competition/api/docs"
	v1 "musobaqa/farm-competition/api/handlers/v1"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"

	"musobaqa/farm-competition/internal/pkg/config"
	tokens "musobaqa/farm-competition/internal/pkg/token"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Product        products.Product
	Animals        animals.Animal
	Food           foods.Food
	Drug           drugs.Drug
	Enforcer       *casbin.Enforcer
}

// NewRoute
// @title Welcome To Farmish API
// @Description API for Farmer
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		JwtHandler:     option.JwtHandler,
		Product:        option.Product,
		Animals:        option.Animals,
		Food:           option.Food,
		Drug:           option.Drug,
		Enforcer:       option.Enforcer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	// router.Use(middleware.Tracing)
	// router.Use(middleware.CheckCasbinPermission(option.Enforcer, *option.Config))

	router.Static("/media", "./media")
	api := router.Group("/v1")

	// ANIMAL METHODS
	api.POST("/animal", HandlerV1.CreateAnimal)
	api.GET("/animals/products/{id}")

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
