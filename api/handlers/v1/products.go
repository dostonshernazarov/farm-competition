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
	"go.opentelemetry.io/otel/attribute"
)

// CREATE PRODUCT
// @Summary CREATE PRODUCT
// @Description Api for Create new product
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param Product body models.ProductReq true "createModel"
// @Success 201 {object} models.ProductRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/products [post]
func (h *HandlerV1) CreateProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.ProductReq
		res  *entity.Product
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

	checkRes, err := h.Product.UniqueProductName(ctx, body.ProductName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	if checkRes != 0 {
		getProduct, err := h.Product.Get(ctx, map[string]string{"name": body.ProductName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}

		res, err = h.Product.Update(ctx, &entity.Product{
			ID:            getProduct.ID,
			Name:          getProduct.Name,
			TotalCapacity: int64(getProduct.TotalCapacity) + int64(body.TotalCapacity),
			Union:         getProduct.Union,
			Description:   getProduct.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}
	} else {

		res, err = h.Product.Create(ctx, &entity.Product{
			ID:            uuid.New().String(),
			Name:          body.ProductName,
			Union:         body.Union,
			TotalCapacity: int64(body.TotalCapacity),
			Description:   body.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, &models.ProductRes{
		Id:            res.ID,
		ProductName:   res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.TotalCapacity),
	})
}

// GET PRODUCT
// @Summary GET PRODUCT BY PRODUCT ID
// @Description Api for Get product by ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/products/{id} [get]
func (h *HandlerV1) GetProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	println("\n", id, "\n")

	res, err := h.Product.Get(ctx, map[string]string{"id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.ProductRes{
		Id:            res.ID,
		ProductName:   res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.TotalCapacity),
	})
}

// LIST PRODUCT
// @Summary LIST PRODUCT
// @Description Api for List Product by page limit and extra values
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.ProductFieldValues true "request"
// @Success 200 {object} models.ListProductsRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/products [get]
func (h *HandlerV1) ListProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListProduct")
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

	name := c.Query("name")
	union := c.Query("union")

	mapP := map[string]interface{}{
		"name":  name,
		"union": union,
	}

	res, err := h.Product.List(ctx, params.Page, params.Limit, mapP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	var resList []*models.ProductRes
	for _, i := range res.Products {
		var resItem models.ProductRes
		resItem.Id = i.ID
		resItem.ProductName = i.Name
		resItem.Description = i.Description
		resItem.Union = i.Union
		resItem.TotalCapacity = int64(i.TotalCapacity)

		resList = append(resList, &resItem)
	}

	c.JSON(http.StatusOK, &models.ListProductsRes{
		Products: resList,
	})
}

// UPDATE
// @Summary UPDATE PRODUCT
// @Description Api for Update product by product id
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param Product body models.ProductRes true "createModel"
// @Success 200 {object} models.ProductRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/products [put]
func (h *HandlerV1) UpdateProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateUser")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.ProductRes
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	res, err := h.Product.Update(ctx, &entity.Product{
		ID:            body.Id,
		Name:          body.ProductName,
		Union:         body.Union,
		TotalCapacity: int64(body.TotalCapacity),
		Description:   body.Description,
	})
	if err != nil {
		c.JSON(500, models.InternalMessage)
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.ProductRes{
		Id:            res.ID,
		ProductName:   res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.TotalCapacity),
	})
}

// DELETE
// @Summary DELETE PRODUCT
// @Description Api for Delete product by product ID
// @Tags PRODUCT
// @Accept json
// @Produce json
// @Param id query string true "ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/products/{id} [delete]
func (h *HandlerV1) DeleteProduct(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteAnimal")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Query("id")

	_, err := h.Product.Get(ctx, map[string]string{"id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("failed to get product in delete", l.Error(err))
		return
	}

	err = h.Product.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete product", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Result{
		Message: "Product has been deleted",
	})
}
