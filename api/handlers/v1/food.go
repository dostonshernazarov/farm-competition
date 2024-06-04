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

// CREATE FOOD
// @Summary CREATE FOOD
// @Description Api for Create new food
// @Tags FOOD
// @Accept json
// @Produce json
// @Param Food body models.FoodReq true "createModel"
// @Success 201 {object} models.FoodRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/foods [post]
func (h *HandlerV1) CreateFood(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateFood")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.FoodReq
		res  *entity.Food
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

	checkRes, err := h.Food.UniqueFoodName(ctx, body.FoodName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	if checkRes != 0 {
		getFood, err := h.Food.Get(ctx, map[string]string{"name": body.FoodName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}

		res, err = h.Food.Update(ctx, &entity.Food{
			ID:          getFood.ID,
			Name:        getFood.Name,
			Capacity:    getFood.Capacity + uint64(body.TotalCapacity),
			Union:       getFood.Union,
			Description: getFood.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}
	} else {

		res, err = h.Food.Create(ctx, &entity.Food{
			ID:          uuid.New().String(),
			Name:        body.FoodName,
			Capacity:    uint64(body.TotalCapacity),
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
	}

	c.JSON(http.StatusCreated, &models.FoodRes{
		Id:            res.ID,
		FoodName:      res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.Capacity),
	})
}

// GET FOOD
// @Summary GET FOOD BY FOOD ID
// @Description Api for Get food by ID
// @Tags FOOD
// @Accept json
// @Produce json
// @Param id path string true "Food ID"
// @Success 200 {object} models.FoodRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/foods/{id} [get]
func (h *HandlerV1) GetFood(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetFood")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	res, err := h.Food.Get(ctx, map[string]string{"id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.FoodRes{
		Id:            res.ID,
		FoodName:      res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.Capacity),
	})
}

// LIST FOOD
// @Summary LIST FOOD
// @Description Api for List food by page limit and extra values
// @Tags FOOD
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.FoodFieldValues true "request"
// @Success 200 {object} models.ListFoodsRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/foods [get]
func (h *HandlerV1) ListFood(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListFood")
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

	mapF := map[string]interface{}{
		"name":  name,
		"union": union,
	}

	res, err := h.Food.List(ctx, params.Page, params.Limit, mapF)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	var resList []*models.FoodRes
	for _, i := range res.Foods {
		var resItem models.FoodRes
		resItem.Id = i.ID
		resItem.FoodName = i.Name
		resItem.Description = i.Description
		resItem.Union = i.Union
		resItem.TotalCapacity = int64(i.Capacity)

		resList = append(resList, &resItem)
	}

	c.JSON(http.StatusOK, &models.ListFoodsRes{
		Foods: resList,
		Count: int64(res.TotalCount),
	})
}

// UPDATE
// @Summary UPDATE FOOD
// @Description Api for Update food by food id
// @Tags FOOD
// @Accept json
// @Produce json
// @Param Food body models.FoodRes true "createModel"
// @Success 200 {object} models.FoodRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/foods [put]
func (h *HandlerV1) UpdateFood(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateFood")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.FoodRes
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	res, err := h.Food.Update(ctx, &entity.Food{
		ID:          body.Id,
		Name:        body.FoodName,
		Capacity:    uint64(body.TotalCapacity),
		Union:       body.Union,
		Description: body.Description,
	})
	if err != nil {
		c.JSON(500, models.InternalMessage)
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.FoodRes{
		Id:            res.ID,
		FoodName:      res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.Capacity),
	})
}

// DELETE
// @Summary DELETE FOOD
// @Description Api for Delete food by food ID
// @Tags FOOD
// @Accept json
// @Produce json
// @Param id path string true "Food ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/foods/{id} [delete]
func (h *HandlerV1) DeleteFood(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteFood")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	_, err := h.Food.Get(ctx, map[string]string{"id": id})
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
		Message: "Food has been deleted",
	})
}
