package v1

import (
	"musobaqa/farm-competition/api/models"
	l "musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"musobaqa/farm-competition/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

// CREATE DELIVERY
// @Summary CREATE DELIVERY
// @Description Api for Create new delivery
// @Tags DELIVERY
// @Accept json
// @Produce json
// @Param Delivery body models.DeliveryReq true "createModel"
// @Success 201 {object} models.DeliveryRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/delivery [post]
func (h *HandlerV1) Createdelivery(c *gin.Context) {
	_, span := otlp.Start(c, "api", "Createdelivery")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.DeliveryReq
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongDateMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	// res, err := h.Delivery.Create(ctx, &entity.Food{
	// 	ID:          uuid.New().String(),
	// 	Name:        body.DeliveryName,
	// 	Capacity:    uint64(body.TotalCapacity),
	// 	Union:       body.Union,
	// 	Description: body.Description,
	// })
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, models.Error{
	// 		Message: models.InternalMessage,
	// 	})
	// 	h.Logger.Error(err.Error())
	// 	return
	// }

	c.JSON(http.StatusCreated, &models.DeliveryRes{
		ID:          "",
		ProductName: "",
		Category:    "",
		Capacity:    0,
		Union:       "",
		Time:        "",
	})
}

// GET DELIVERY
// @Summary GET DELIVERY BY FOOD ID
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

	_, err := h.Food.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.DeliveryRes{
		ID:          id,
		ProductName: "",
		Category:    "",
		Capacity:    0,
		Union:       "",
		Time:        "",
	})
}

// LIST DELIVERY
// @Summary LIST DELIVERY
// @Description Api for List delivery by page limit and extra values
// @Tags DELIVERY
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.FoodFieldValues true "request"
// @Success 200 {object} models.ListDeliverysRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/delivery [get]
func (h *HandlerV1) ListDelivery(c *gin.Context) {
	_, span := otlp.Start(c, "api", "ListDelivery")
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

	_ = map[string]interface{}{
		"name":  name,
		"union": union,
	}

	// res, err := h.Delivery.List(ctx, params.Page, params.Limit)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, models.Error{
	// 		Message: models.InternalMessage,
	// 	})
	// 	h.Logger.Error(err.Error())
	// 	return
	// }

	// var resList []*models.DeliveryRes
	// for _, i := range res.Deliverys {
	// 	var resItem models.DeliveryRes
	// 	resItem.Id = i.ID
	// 	resItem.DeliveryName = i.Name
	// 	resItem.Description = i.Description
	// 	resItem.Union = i.Union
	// 	resItem.TotalCapacity = int64(i.Capacity)

	// 	resList = append(resList, &resItem)
	// }

	c.JSON(http.StatusOK, &models.ListDeliverysRes{
		Delivery: []*models.DeliveryRes{},
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
// @Router /v1/deliverys [put]
func (h *HandlerV1) UpdateDelivery(c *gin.Context) {
	_, span := otlp.Start(c, "api", "UpdateDelivery")
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

	// res, err := h.Delivery.Update(ctx, &entity.Delivery{
	// 	ID:          body.Id,
	// 	Name:        body.DeliveryName,
	// 	Capacity:    uint64(body.TotalCapacity),
	// 	Union:       body.Union,
	// 	Description: body.Description,
	// })
	// if err != nil {
	// 	c.JSON(500, models.InternalMessage)
	// 	h.Logger.Error(err.Error())
	// 	return
	// }

	c.JSON(http.StatusOK, &models.DeliveryRes{

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

	_, err := h.Food.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error(err.Error())
		return
	}

	err = h.Food.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete product", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Result{
		Message: "Delivery has been deleted",
	})
}