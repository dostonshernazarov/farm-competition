package v1

import (
	"musobaqa/farm-competition/api/models"
	"musobaqa/farm-competition/internal/entity"
	l "musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"musobaqa/farm-competition/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
)

// CREATE ANIMAL
// @Summary CREATE ANIMAL
// @Description Api for Create animal
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param Animal body models.AnimalReq true "createModel"
// @Success 201 {object} models.AnimalRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animal [post]
func (h *HandlerV1) CreateAnimal(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateAnimal")
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
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongDateMessage,
		})
		l.Error(err)
		return
	}

	err = body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error(err.Error())
		return
	}

	res, err := h.Animals.Create(ctx, &entity.Animal{
		ID:          uuid.NewString(),
		Name:        body.Name,
		CategoryName:  body.CategoryName,
		Gender:      body.Gender,
		BirthDay:    body.DateOfBirth,
		Genus:       body.Genus,
		Weight:      uint64(body.Weight),
		IsHealth:    cast.ToString(body.IsHealth),
		Description: body.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, &models.AnimalRes{
		Id:           res.ID,
		Name:         res.Name,
		CategoryName: res.CategoryName,
		Gender:       res.Gender,
		DateOfBirth:  res.BirthDay,
		Description:  res.Description,
		Genus:        res.Genus,
		Weight:       float32(res.Weight),
		IsHealth:     cast.ToBool(res.IsHealth),
	})
}

// GET ANIMAL
// @Summary GET ANIMAL BY ANIMAL ID
// @Description Api for Get Animal by ID
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param id path string true "Animal_id"
// @Success 200 {object} models.AnimalRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/{id} [get]
func (h *HandlerV1) GetAnimal(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetAnimalWithProdactList	")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	res, err := h.Animals.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		l.Error(err)
		return
	}

	c.JSON(http.StatusOK, &models.AnimalRes{
		Id:           res.ID,
		Name:         res.Name,
		CategoryName: res.CategoryName,
		Gender:       res.Gender,
		DateOfBirth:  res.BirthDay,
		Description:  res.Description,
		Genus:        res.Genus,
		Weight:       float32(res.Weight),
		IsHealth:     cast.ToBool(res.IsHealth),
	})
}

// GET ANIMAL
// @Summary GET ANIMAL BY ANIMAL ID WITH PRODUCTS
// @Description Api for Get Animal with products
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param animal_id path string true "ID"
// @Param request query models.Pagination true "request"
// @Success 200 {object} models.AnimalProdactList
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
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
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongDateMessage,
		})
		return
	}

	println(params)
	c.JSON(http.StatusOK, &models.AnimalProdactList{})
}

// GET ANIMAL
// @Summary GET ANIMAL BY ANIMAL ID WITH FOODS
// @Description Api for Get Animal by animal id with foods
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param animal_id path string true "ID"
// @Param request query models.Pagination true "request"
// @Success 200 {object} models.AnimalProdactList
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
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
		c.JSON(http.StatusBadRequest, models.WrongInfoMessage)
		return
	}

	println(params)
	c.JSON(http.StatusOK, &models.AnimalFoodList{})
}

// LIST ANIMALS
// @Summary LIST ANIMALS
// @Description Api for List Animals by page limit and extra values
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.AnimalFieldValues true "request"
// @Success 200 {object} models.ListAnimalsRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals [get]
func (h *HandlerV1) ListAnimals(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListAnimals")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		return
	}

	println(params)

	category := c.Query("category")
	genus := c.Query("genus")
	gender := c.Query("gender")
	weight := c.Query("weight")
	is_health := c.Query("is_health")

	// if strings.ToLower(gender) != "male" || strings.ToLower(gender) != "female" {
	// 	gender = ""
	// }

	mapA := map[string]interface{}{
		"category":  category,
		"genus":     genus,
		"gender":    gender, // Empty string
		"weight":    weight,
		"is_health": is_health,
	}

	list, err := h.Animals.List(ctx, params.Page, params.Limit, mapA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error("failed to get list animal", l.Error(err))
		return
	}

	var reslist []*models.AnimalRes
	for _, i := range list.Animals {
		var res models.AnimalRes
		res.Id = i.ID
		res.CategoryName = i.Name
		res.DateOfBirth = i.BirthDay
		res.Description = i.Description
		res.Gender = i.Gender
		res.Genus = i.Genus
		res.Weight = float32(i.Weight)
		res.IsHealth = cast.ToBool(i.IsHealth)
		reslist = append(reslist, &res)
	}

	c.JSON(http.StatusOK, &models.ListAnimalsRes{
		Animals: reslist,
		Count: int64(list.TotalCount),
	})
}

// UPDATE
// @Summary UPDATE ANIMAL
// @Description Api for Update Animal
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param Animal body models.AnimalRes true "createModel"
// @Success 200 {object} models.AnimalRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals [put]
func (h *HandlerV1) UpdateAnimal(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimalRes
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	resAnimals, err := h.Animals.Update(ctx, &entity.Animal{
		ID:          body.Id,
		Name:        body.Name,
		CategoryName:  body.CategoryName,
		Gender:      body.Gender,
		BirthDay:    body.DateOfBirth,
		Genus:       body.Genus,
		Weight:      uint64(body.Weight),
		IsHealth:    cast.ToString(body.IsHealth),
		Description: body.Description,
	})
	if err != nil {
		c.JSON(500, models.InternalMessage)
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.AnimalRes{
		Id:           resAnimals.ID,
		Name:         resAnimals.Name,
		CategoryName: resAnimals.CategoryName,
		Gender:       resAnimals.Gender,
		DateOfBirth:  resAnimals.BirthDay,
		Description:  resAnimals.Description,
		Genus:        resAnimals.Genus,
		Weight:       float32(resAnimals.Weight),
		IsHealth:     cast.ToBool(resAnimals.IsHealth),
	})
}

// DELETE
// @Summary DELETE ANIMAL
// @Description Api for Delete animal by animal ID
// @Tags ANIMAL
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/{id} [delete]
func (h *HandlerV1) DeleteAnimal(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteAnimal")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	_, err := h.Animals.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("failed to get animal in delete", l.Error(err))
		return
	}

	err = h.Animals.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete animal", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Result{
		Message: "Animal has been deleted",
	})
}

// LIST CATEGORY
// @Summary LIST ANIMAL CATEGORY
// @Description Api for List AnimalCategory
// @Tags CATEGORY
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Success 200 {object} models.CategoryRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/category [get]
func (h *HandlerV1) ListCategory(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListAnimals")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParam(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		return
	}

	println(params, ctx)
}
