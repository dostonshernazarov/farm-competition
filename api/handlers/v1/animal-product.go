package v1

import (
	"musobaqa/farm-competition/api/models"
	"musobaqa/farm-competition/internal/entity"
	l "musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"musobaqa/farm-competition/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
)

// CREATE ANIMAL PRODUCT
// @Summary CREATE ANIMAL PRODUCT
// @Description Api for Create product which has got from animal
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param Animal-Product body models.AnimalProductReq true "createModel"
// @Success 201 {object} models.AnimalProductRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/products [post]
func (h *HandlerV1) CreateAnimalProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimalProductReq
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongDateMessage,
		})
		h.Logger.Error(err.Error())
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

	res, err := h.AnimalProduct.Create(ctx, &entity.AnimalProductReq{
		AnimalID:  body.AnimalID,
		ProductID: body.ProductID,
		Capacity:  body.Capacity,
		GetTime:   body.GetTime,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	_, err = h.Product.Update(ctx, &entity.Product{
		ID:            body.ProductID,
		Name:          res.Product.Name,
		Union:         res.Product.Union,
		TotalCapacity: res.Product.TotalCapacity + body.Capacity,
		Description:   res.Product.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, &models.AnimalProductRes{
		Id:             res.ID,
		AnimalID:       res.Animal.ID,
		AnimalName:     res.Animal.Name,
		AnimalCategory: res.Animal.CategoryName,
		ProductName:    res.Product.Name,
		Capacity:       res.Capacity,
		Union:          res.Product.Union,
		GetTime:        res.GetTime,
	})
}

// GET ANIMAL PRODUCT
// @Summary GET ANIMAL PRODUCT BY ID
// @Description Api for Get product which has got from animal by ID
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param id path string true "Animal Product ID"
// @Success 200 {object} models.AnimalProductRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/products/{id} [get]
func (h *HandlerV1) GetAnimalProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	res, err := h.AnimalProduct.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.AnimalProductRes{
		Id:             res.ID,
		AnimalID:       res.Animal.ID,
		AnimalName:     res.Animal.Name,
		AnimalCategory: res.Animal.CategoryName,
		ProductName:    res.Product.Name,
		Capacity:       res.Capacity,
		Union:          res.Product.Union,
		GetTime:        res.GetTime,
	})
}

// LIST ANIMAL PRODUCT
// @Summary LIST ANIMAL PRODUCT
// @Description Api for List products which have got from animals by page limit and extra values
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.AnimalProductFieldValues true "request"
// @Success 200 {object} models.ListAnimalProductsRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/products [get]
func (h *HandlerV1) ListAnimalProducts(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListAnimalProducts")
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

	get_time := c.Query("get_time")

	mapA := map[string]interface{}{
		"get_time": get_time,
	}

	res, err := h.AnimalProduct.List(ctx, params.Page, params.Limit, mapA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	var resList []*models.AnimalProductRes
	for _, i := range res.AnimalProducts {
		var resItem models.AnimalProductRes
		resItem.Id = i.ID
		resItem.AnimalID = i.Animal.ID
		resItem.AnimalName = i.Animal.Name
		resItem.Union = i.Product.Union
		resItem.AnimalCategory = i.Animal.CategoryName
		resItem.ProductName = i.Product.Name
		resItem.Capacity = i.Capacity
		resItem.Union = i.Product.Union
		resItem.GetTime = i.GetTime

		resList = append(resList, &resItem)
	}

	c.JSON(http.StatusOK, &models.ListAnimalProductsRes{
		AnimalProducts: resList,
		Count:          int64(res.TotalCount),
	})
}

// UPDATE ANIMAL PRODUCT
// @Summary UPDATE ANIMAL PRODUCT
// @Description Api for Update product which has got from animal
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param Animal-Product body models.AnimalProductUpdateReq true "createModel"
// @Success 200 {object} models.AnimalProductRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/products [put]
func (h *HandlerV1) UpdateAnimalProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimalProductUpdateReq
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	res, err := h.AnimalProduct.Update(ctx, &entity.AnimalProductReq{
		ID:        body.ID,
		AnimalID:  body.AnimalID,
		ProductID: body.ProductID,
		Capacity:  body.Capacity,
		GetTime:   body.GetTime,
	})
	if err != nil {
		c.JSON(500, models.InternalMessage)
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.AnimalProductRes{
		Id:             res.ID,
		AnimalID:       res.Animal.ID,
		AnimalName:     res.Animal.Name,
		AnimalCategory: res.Animal.CategoryName,
		ProductName:    res.Product.Name,
		Capacity:       res.Capacity,
		Union:          res.Product.Union,
		GetTime:        res.GetTime,
	})
}

// DELETE ANIMAL PRODUCT
// @Summary DELETE ANIMAL PRODUCT
// @Description Api for Delete product which has got from animal by ID
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param id path string true "Animal Product ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/products/{id} [delete]
func (h *HandlerV1) DeleteAnimalProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	_, err := h.AnimalProduct.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("failed to get animal product in delete", l.Error(err))
		return
	}

	err = h.AnimalProduct.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete animal product", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Result{
		Message: "Animal product has been deleted",
	})
}


// LIST ANIMAL PRODUCTS BY ANIMAL ID
// @Summary LIST ANIMAL PRODUCTS BY ANIMAL ID
// @Description Api for List products with animal which have got from animals by animal ID
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.AnimalProductByAnimalIdFieldValues true "request"
// @Success 200 {object} models.AnimalProductByAnimalIdRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animal-products [get]
func (h *HandlerV1) ListAnimalProductsByAnimalID(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListAnimalProductsByAnimalID")
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

	animal_id := c.Query("animal_id")
	if animal_id == "" {
		c.JSON(http.StatusBadGateway, models.WrongInfoMessage)
		h.Logger.Error("empty animal id",)
		return
	}

	res, err := h.AnimalProduct.ListProducts(ctx, params.Page, params.Limit, animal_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to list animal product by product id", l.Error(err))
		return
	}

	var resList []*models.ProductRes
	for _, i := range res.Products {
		var resItem models.ProductRes
		resItem.Id = i.ID
		resItem.ProductName = i.Name
		resItem.Union = i.Union
		resItem.Description = i.Description
		resItem.TotalCapacity = i.TotalCapacity
		resList = append(resList, &resItem)
	}

	c.JSON(http.StatusOK, &models.AnimalProductByAnimalIdRes{
		Animal: &models.AnimalRes{
			Id:           res.Animal.ID,
			Name:         res.Animal.Name,
			CategoryName: res.Animal.CategoryName,
			Gender:       res.Animal.Gender,
			DateOfBirth:  res.Animal.BirthDay,
			Description:  res.Animal.Description,
			Genus:        res.Animal.Genus,
			Weight:       float32(res.Animal.Weight),
			IsHealth:     cast.ToBool(res.Animal.IsHealth),
		},
		Products: resList,
		Count:          int64(res.TotalCount),
	})
}

// LIST PRODUCT ANIMALS BY PRODUCT ID
// @Summary LIST  PRODUCT ANIMALS BY PRODUCT ID
// @Description Api for List animals with product which have got from animals by product ID
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.AnimalProductByProductIdFieldValues true "request"
// @Success 200 {object} models.AnimalProductByProductIdRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/product-animals [get]
func (h *HandlerV1) ListAnimalProductsByProductID(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListAnimalProductsByProductID")
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

	product_id := c.Query("product_id")
	if product_id == "" {
		c.JSON(http.StatusBadGateway, models.WrongInfoMessage)
		h.Logger.Error("empty product id",)
		return
	}

	res, err := h.AnimalProduct.ListAnimals(ctx, params.Page, params.Limit, product_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to list animal product by animal id", l.Error(err))
		return
	}


	var resList []*models.AnimalCapRes
	for _, i := range res.Animals {
		var resItem models.AnimalCapRes
		resItem.Id = i.ID
		resItem.Name = i.Name
		resItem.CategoryName = i.CategoryName
		resItem.Gender = i.Gender
		resItem.DateOfBirth = i.BirthDay
		resItem.Description = i.Description
		resItem.Genus = i.Genus
		resItem.Weight = float32(i.Weight)
		resItem.IsHealth = cast.ToBool(i.IsHealth)
		resItem.TotalCapacity = i.TotalCapacity

		resList = append(resList, &resItem)
	}

	c.JSON(http.StatusOK, &models.AnimalProductByProductIdRes{
		Product: models.ProductRes{
			Id:            res.Product.ID,
			ProductName:   res.Product.Name,
			Union:         res.Product.Union,
			Description:   res.Product.Description,
			TotalCapacity: res.Product.TotalCapacity,
		},
		Animals: resList,
		Count:          int64(res.TotalCount),
	})
}
