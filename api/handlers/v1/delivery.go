package v1

import (
	"errors"
	"musobaqa/farm-competition/api/models"
	"musobaqa/farm-competition/internal/entity"
	l "musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"musobaqa/farm-competition/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.opentelemetry.io/otel/attribute"
)

// CREATE DELIVERY
// @Summary CREATE DELIVERY
// @Description Api for Create new delivery
// @Tags DELIVERY
// @Accept json
// @Produce json
// @Param Delivery body models.DeliveryCreateReq true "createModel"
// @Success 201 {object} models.DeliveryCreateRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/delivery [post]
func (h *HandlerV1) CreateDelivery(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "Createdelivery")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.DeliveryCreateReq
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

	_, err = h.Delivery.Create(ctx, &entity.Delivery{
		Name:     body.ProductName,
		Category: body.Category,
		Capacity: body.Capacity,
		Union:    body.Union,
		Time:     body.Time,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	if body.Category == "food" {
		foodRes, err := h.Food.Get(ctx, map[string]string{"name": body.ProductName})

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				println("\n check \n")
				_, err := h.Food.Create(ctx, &entity.Food{
					ID:          uuid.New().String(),
					Name:        body.ProductName,
					Capacity:    uint64(body.Capacity),
					Union:       body.Union,
					Description: body.Description,
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, models.Error{
						Message: models.InternalMessage,
					})
					h.Logger.Error(err.Error())
					return
				}
				c.JSON(http.StatusCreated, &models.DeliveryCreateRes{
					Message: "Product successfully created",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}

		_, err = h.Food.Update(ctx, &entity.Food{
			ID:          foodRes.ID,
			Name:        foodRes.Name,
			Capacity:    foodRes.Capacity + uint64(body.Capacity),
			Union:       foodRes.Union,
			Description: foodRes.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}
		c.JSON(http.StatusCreated, &models.DeliveryCreateRes{
			Message: "Product successfully updated",
		})
		return

	}

	if body.Category == "drug" {
		drugRes, err := h.Drug.Get(ctx, map[string]string{"name": body.ProductName})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				_, err := h.Drug.Create(ctx, &entity.Drug{
					ID:          uuid.NewString(),
					Name:        body.ProductName,
					Status:      body.Status,
					Capacity:    uint64(body.Capacity),
					Union:       body.Union,
					Description: body.Description,
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, models.Error{
						Message: models.InternalMessage,
					})
					h.Logger.Error(err.Error())
					return
				}
				c.JSON(http.StatusCreated, &models.DeliveryCreateRes{
					Message: "Product successfully created",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}

		_, err = h.Drug.Update(ctx, &entity.Drug{
			ID:          drugRes.ID,
			Name:        drugRes.Name,
			Status:      drugRes.Status,
			Capacity:    drugRes.Capacity + uint64(body.Capacity),
			Union:       drugRes.Union,
			Description: drugRes.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}
		c.JSON(http.StatusCreated, &models.DeliveryCreateRes{
			Message: "Product successfully updated",
		})
		return
	}

	c.JSON(http.StatusCreated, &models.DeliveryCreateRes{
		Message: "Warning! Went wrong",
	})
}

// GET DELIVERY
// @Summary GET DELIVERY BY ID
// @Description Api for Get delivery by ID
// @Tags DELIVERY
// @Accept json
// @Produce json
// @Param id path string true "Delivery ID"
// @Success 200 {object} models.DeliveryRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/delivery/{id} [get]
func (h *HandlerV1) GetDelivery(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetDelivery")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	res, err := h.Delivery.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.DeliveryRes{
		ID:          res.ID,
		ProductName: res.Name,
		Category:    res.Category,
		Capacity:    res.Capacity,
		Union:       res.Union,
		Time:        res.Time,
	})
}

// LIST DELIVERY
// @Summary LIST DELIVERY
// @Description Api for List delivery by page limit and extra values
// @Tags DELIVERY
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.DeliveryFieldValues true "request"
// @Success 200 {object} models.ListDeliverysRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/delivery [get]
func (h *HandlerV1) ListDelivery(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListDelivery")
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

	name := c.Query("name")
	category := c.Query("category")
	time := c.Query("time")

	mapD := map[string]interface{}{
		"name":     name,
		"category": category,
		"time":     time,
	}

	res, err := h.Delivery.List(ctx, params.Page, params.Limit, mapD)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	var resList []*models.DeliveryRes
	for _, i := range res.Deliveries {
		var resItem models.DeliveryRes
		resItem.ID = i.ID
		resItem.ProductName = i.Name
		resItem.Category = i.Category
		resItem.Union = i.Union
		resItem.Capacity = int64(i.Capacity)
		resItem.Time = i.Time

		resList = append(resList, &resItem)
	}

	c.JSON(http.StatusOK, &models.ListDeliverysRes{
		Delivery: resList,
		Count:    res.TotalCount,
	})
}

// UPDATE
// @Summary UPDATE DELIVERY
// @Description Api for Update delivery by food id
// @Tags DELIVERY
// @Accept json
// @Produce json
// @Param Delivery body models.DeliveryRes true "createModel"
// @Success 200 {object} models.DeliveryRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/delivery [put]
func (h *HandlerV1) UpdateDelivery(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateDelivery")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.DeliveryRes
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	res, err := h.Delivery.Update(ctx, &entity.Delivery{
		ID:       body.ID,
		Name:     body.ProductName,
		Capacity: body.Capacity,
		Union:    body.Union,
		Time:     body.Time,
	})
	if err != nil {
		c.JSON(500, models.InternalMessage)
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.DeliveryRes{
		ID:          res.ID,
		ProductName: res.Name,
		Category:    res.Category,
		Capacity:    res.Capacity,
		Union:       res.Union,
		Time:        res.Time,
	})
}

// DELETE
// @Summary DELETE DELIVERY
// @Description Api for Delete delivery by delivery ID
// @Tags DELIVERY
// @Accept json
// @Produce json
// @Param id path string true "Delivery ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/delivery/{id} [delete]
func (h *HandlerV1) DeleteDelivery(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteDelivery")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	_, err := h.Delivery.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error(err.Error())
		return
	}

	err = h.Delivery.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete product", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Result{
		Message: "Delivery has been deleted",
	})
}
