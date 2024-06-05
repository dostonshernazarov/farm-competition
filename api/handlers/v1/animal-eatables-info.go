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

	var dailyReq []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	}

	for _, value := range body.Daily {
		dailyReq = append(dailyReq, struct {
			Capacity int64  "json:\"capacity\""
			Time     string "json:\"time\""
		}{
			Time:     value.Time,
			Capacity: value.Capacity,
		})
	}

	eatablesRes, err := h.EatablesInfo.Create(ctx, &entity.Eatables{
		AnimalID:  body.AnimalID,
		EatableID: body.EatablesID,
		Category:  body.Category,
		Daily:     dailyReq,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	var res []*models.Daily
	for _, value := range eatablesRes.Daily {
		res = append(res, &models.Daily{
			Time:     value.Time,
			Capacity: value.Capacity,
		})
	}

	c.JSON(http.StatusCreated, &models.AnimaEatablesInfoRes{
		ID:         eatablesRes.ID,
		AnimalID:   eatablesRes.AnimalID,
		EatablesID: eatablesRes.Eatable.ID,
		Daily:      res,
		Category:   eatablesRes.Category,
	})
}

// UPDATE ANIMAL EATABLES INFO
// @Summary UPDATE ANIMAL EATABLES INFO
// @Description Api for Update animal eatables info
// @Tags EATABLES-INFO
// @Accept json
// @Produce json
// @Param Animal-Eatables body models.AnimaEatablesInfoRes true "UpdateModel"
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
		body models.AnimaEatablesInfoRes
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	var dailyReq []struct {
		Capacity int64  `json:"capacity"`
		Time     string `json:"time"`
	}

	for _, value := range body.Daily {
		dailyReq = append(dailyReq, struct {
			Capacity int64  `json:"capacity"`
			Time     string `json:"time"`
		}{
			Capacity: value.Capacity,
			Time:     value.Time,
		})
	}

	res, err := h.EatablesInfo.Update(ctx, &entity.Eatables{
		ID:        body.ID,
		AnimalID:  body.AnimalID,
		EatableID: body.EatablesID,
		Category:  body.Category,
		Daily:     dailyReq,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error(err.Error())
		return
	}

	var resDaily []*models.Daily
	for _, value := range res.Daily {
		resDaily = append(resDaily, &models.Daily{
			Time:     value.Time,
			Capacity: value.Capacity,
		})
	}

	c.JSON(http.StatusOK, &models.AnimaEatablesInfoRes{
		ID:         res.ID,
		AnimalID:   res.AnimalID,
		EatablesID: res.Eatable.ID,
		Daily:      resDaily,
		Category:   res.Category,
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

	err := h.EatablesInfo.Delete(ctx, id)
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
// @Success 200 {object} models.ListFootInfoByAnimalRes
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

	animalID := c.Query("animal_id")
	if animalID == "" {
		c.JSON(http.StatusBadGateway, models.WrongInfoMessage)
		h.Logger.Error("empty animal id")
		return
	}

	res, err := h.EatablesInfo.GetFoods(ctx, params.Page, params.Limit, animalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to list animal eatables info by animal id", l.Error(err))
		return
	}

	var resList models.ListFootInfoByAnimalRes
	for _, i := range res.Eatables {
		var (
			dailyInfo []*models.Daily
			resItem   models.AnimaFoodInfoRes
		)

		for _, value := range i.Daily {
			dailyInfo = append(dailyInfo, &models.Daily{
				Time:     value.Time,
				Capacity: value.Capacity,
			})
		}
		resItem.ID = i.ID
		resItem.Eatables.Id = i.Food.ID
		resItem.Eatables.FoodName = i.Food.Name
		resItem.Eatables.Union = i.Food.Union
		resItem.Eatables.Description = i.Food.Description
		resItem.Eatables.TotalCapacity = int64(i.Food.Capacity)
		resItem.AnimalID = i.AnimalID
		resItem.Category = "food"
		resItem.Daily = dailyInfo
		resList.Eatables = append(resList.Eatables, &resItem)
	}
	resList.Count = int64(res.TotalCount)

	c.JSON(http.StatusOK, &models.ListFootInfoByAnimalRes{
		Eatables: resList.Eatables,
		Count:    resList.Count,
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

	animalID := c.Query("animal_id")
	if animalID == "" {
		c.JSON(http.StatusBadGateway, models.WrongInfoMessage)
		h.Logger.Error("empty animal id")
		return
	}

	res, err := h.EatablesInfo.GetDrugs(ctx, params.Page, params.Limit, animalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to list animal eatables info by animal id", l.Error(err))
		return
	}

	var resList models.ListDrugInfoByAnimalRes
	for _, i := range res.Eatables {
		var (
			dailyInfo []*models.Daily
			resItem   models.AnimalDrugInfoRes
		)

		for _, value := range i.Daily {
			dailyInfo = append(dailyInfo, &models.Daily{
				Time:     value.Time,
				Capacity: value.Capacity,
			})
		}
		resItem.ID = i.ID
		resItem.Eatables.Id = i.Drug.ID
		resItem.Eatables.DrugName = i.Drug.Name
		resItem.Eatables.Status = i.Drug.Status
		resItem.Eatables.Union = i.Drug.Union
		resItem.Eatables.Description = i.Drug.Description
		resItem.Eatables.TotalCapacity = int64(i.Drug.Capacity)
		resItem.AnimalID = i.AnimalID
		resItem.Category = "drug"
		resItem.Daily = dailyInfo
		resList.Eatables = append(resList.Eatables, &resItem)
	}
	resList.Count = int64(res.TotalCount)

	c.JSON(http.StatusOK, &models.ListDrugInfoByAnimalRes{
		Eatables: resList.Eatables,
		Count:    resList.Count,
	})
}
