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

// CREATE ANIMAL EATABLES INFO
// @Summary CREATE ANIMAL EATABLES INFO
// @Description Api for Create animal eatables info
// @Tags EATABLES-INFO
// @Accept json
// @Produce json
// @Param Animal-Eatables body models.AnimaEatablesInfoReq true "createModel"
// @Success 201 {object} models.AnimaEatablesInfoRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/eatables [post]
func (h *HandlerV1) CreateEatablesInfo(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimaEatablesInfoReq
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

	c.JSON(http.StatusCreated, &models.AnimaEatablesInfoRes{
		ID:         "",
		AnimalID:   "",
		EatablesID: "",
		Daily:      []*models.Daily{},
		Category:   "",
	})
}

// UPDATE ANIMAL EATABLES INFO
// @Summary UPDATE ANIMAL EATABLES INFO
// @Description Api for Update animal eatables info
// @Tags EATABLES-INFO
// @Accept json
// @Produce json
// @Param Animal-Eatables body models.AnimaEatablesInfoReq true "createModel"
// @Success 200 {object} models.AnimaEatablesInfoRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/eatables [put]
func (h *HandlerV1) UpdateEatablesInfo(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateAnimalEatablesInfo")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimaEatablesInfoReq
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

	c.JSON(http.StatusOK, &models.AnimaEatablesInfoRes{
		ID:         "",
		AnimalID:   "",
		EatablesID: "",
		Daily:      []*models.Daily{},
		Category:   "",
	})
}

// DELETE ANIMAL EATABLES INFO
// @Summary DELETE ANIMAL EATABLES INFO
// @Description Api for Delete animal eatables info
// @Tags EATABLES-INFO
// @Accept json
// @Produce json
// @Param id path string true "Animal Eatables ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/eatables/{id} [delete]
func (h *HandlerV1) DeleteEatablesInfo(c *gin.Context) {
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
		Message: "Animal eatables info has been deleted",
	})
}

// LIST ANIMAL FOOD INFO ANIMAL ID
// @Summary LIST ANIMAL FOOD INFO ANIMAL ID
// @Description Api for List animal food info by animal ID
// @Tags EATABLES-INFO
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.ListEatablesInfoByAnimalReq true "request"
// @Success 200 {object} models.ListFoodInfoByAnimalRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/food-info [get]
func (h *HandlerV1) ListFoodInfoByAnimalID(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListAnimalEatablesInfoByAnimalID")
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
		h.Logger.Error("empty animal id")
		return
	}

	res, err := h.AnimalProduct.ListProducts(ctx, params.Page, params.Limit, animal_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to list animal eatables info by animal id", l.Error(err))
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

	c.JSON(http.StatusOK, &models.ListFootInfoByAnimalRes{
		Eatables: []*models.AnimaFoodInfoRes{},
		Count: 21,
	})
}


// LIST ANIMAL DRUG INFO ANIMAL ID
// @Summary LIST ANIMAL DRUG INFO ANIMAL ID
// @Description Api for List animal drug info by animal ID
// @Tags EATABLES-INFO
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.ListEatablesInfoByAnimalReq true "request"
// @Success 200 {object} models.ListDrugInfoByAnimalRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/drug-info [get]
func (h *HandlerV1) ListDrugInfoByAnimalID(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListAnimalEatablesInfoByAnimalID")
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
		h.Logger.Error("empty animal id")
		return
	}

	res, err := h.AnimalProduct.ListProducts(ctx, params.Page, params.Limit, animal_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to list animal eatables info by animal id", l.Error(err))
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

	c.JSON(http.StatusOK, &models.ListDrugInfoByAnimalRes{
		Eatables: []*models.AnimaDrugInfoRes{},
		Count: 21,
	})
}
