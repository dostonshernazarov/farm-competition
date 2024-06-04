package v1

import (
	"musobaqa/farm-competition/api/models"
	"musobaqa/farm-competition/internal/entity"
	l "musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"musobaqa/farm-competition/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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
