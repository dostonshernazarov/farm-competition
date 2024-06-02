package v1

import (
	"musobaqa/farm-competition/api/models"
	l "musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"musobaqa/farm-competition/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

// CREATE ANIMAL
// @Summary CREATE ANIMAL
// @Description Api for Create animal
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param Animal body models.AnimalReq true "createModel"
// @Success 201 {object} models.AnimalCreateRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/animal [post]
func (h *HandlerV1) CreateAnimal(c *gin.Context) {
	_, span := otlp.Start(c, "api", "CreateAnimal")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimalReq
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error binding JSON",
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusCreated, &models.AnimalCreateRes{
		Id:           uuid.New().String(),
		Name:         body.Name,
		CategoryName: body.CategoryName,
		Gender:       body.Gender,
		DateOfBirth:  body.DateOfBirth,
		Description:  body.Description,
		Genus:        body.Genus,
		Weight:       body.Weight,
		IsHealth:     body.IsHealth,
	})
}

// GET ANIMAL
// @Summary GET ANIMAL WITH PRODUCTS
// @Description Api for Get Animal with products
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param animal_id path string true "ID"
// @Param request query models.Pagination true "request"
// @Success 200 {object} models.AnimalProdactList
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/animals/product/{id} [get]
func (h *HandlerV1) GetAnimalWithProducts(c *gin.Context) {
	_, span := otlp.Start(c, "api", "GetAnimalWithProdactList	")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		return
	}

	println(params)
	c.JSON(http.StatusOK, &models.AnimalProdactList{})
}

// GET ANIMAL
// @Summary GET ANIMAL WITH FOODS
// @Description Api for Get Animal with foods
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param animal_id path string true "ID"
// @Param request query models.Pagination true "request"
// @Success 200 {object} models.AnimalProdactList
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/animals/food/{id} [get]
func (h *HandlerV1) GetAnimalWithEatables(c *gin.Context) {
	_, span := otlp.Start(c, "api", "GetAnimalWithFoodList	")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		return
	}

	println(params)
	c.JSON(http.StatusOK, &models.AnimalFoodList{})
}

// LIST ANIMALS
// @Summary LIST ANIMALS
// @Security BearerAuth
// @Description Api for ListAnimals
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.AnimalFieldValues true "request"
// @Success 200 {object} models.ListAnimalsRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/animals [get]
func (h *HandlerV1) ListAnimals(c *gin.Context) {
	_, span := otlp.Start(c, "api", "ListUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()


	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		return
	}

	println(params)

	c.JSON(http.StatusOK, &models.ListAnimalsRes{})
}
