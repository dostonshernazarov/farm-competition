package api

import (
	animalproduct "musobaqa/farm-competition/internal/usecase/animal-product"
	"musobaqa/farm-competition/internal/usecase/animals"
	"musobaqa/farm-competition/internal/usecase/delivery"
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
	Delivery delivery.Delivery
	AnimalProduct animalproduct.AnimalProduct
	Enforcer       *casbin.Enforcer
}

// NewRoute
// @title Welcome To Farmish API
// @Description API for Farmer
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
		Delivery: option.Delivery,
		AnimalProduct: option.AnimalProduct,
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
	api.GET("/animals/:id", HandlerV1.GetAnimal)
	api.GET("/animals", HandlerV1.ListAnimals)
	api.PUT("/animals", HandlerV1.UpdateAnimal)
	api.DELETE("/animals/:id", HandlerV1.DeleteAnimal)

	// PRODUCT METHODS
	api.POST("/products", HandlerV1.CreateProduct)
	api.GET("/products/:id", HandlerV1.GetProduct)
	api.GET("/products", HandlerV1.ListProduct)
	api.PUT("/products", HandlerV1.UpdateProduct)
	api.DELETE("/products/:id", HandlerV1.DeleteProduct)

	// DRUG METHODS
	api.POST("/drugs", HandlerV1.CreateDrug)
	api.GET("/drugs/:id", HandlerV1.GetDrug)
	api.GET("/drugs", HandlerV1.ListDrug)
	api.PUT("/drugs", HandlerV1.UpdateDrug)
	api.DELETE("/drugs/:id", HandlerV1.DeleteDrug)

	// FOOD METHODS
	api.POST("/foods", HandlerV1.CreateFood)
	api.GET("/foods/:id", HandlerV1.GetFood)
	api.GET("/foods", HandlerV1.ListFood)
	api.PUT("/foods", HandlerV1.UpdateFood)
	api.DELETE("/foods/:id", HandlerV1.DeleteFood)

	// DELIVERY METHODS
	api.POST("/delivery", HandlerV1.CreateDelivery)
	api.GET("/delivery/:id", HandlerV1.GetDelivery)
	api.GET("/delivery", HandlerV1.ListDelivery)
	api.PUT("/delivery", HandlerV1.UpdateDelivery)
	api.DELETE("/delivery/:id", HandlerV1.DeleteDelivery)

	// ANIMAL PRODUCT METHODS
	api.POST("/animals/products", HandlerV1.CreateAnimalProduct)
	api.GET("/animals/products/:id", HandlerV1.GetAnimalProduct)
	api.GET("/animals/products", HandlerV1.ListAnimalProducts)
	api.PUT("/animals/products", HandlerV1.UpdateAnimalProduct)
	api.DELETE("/animals/products/:id", HandlerV1.DeleteAnimalProduct)
	api.GET("/animal-products", HandlerV1.ListAnimalProductsByAnimalID)
	api.GET("/product-animals", HandlerV1.ListAnimalProductsByProductID)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
