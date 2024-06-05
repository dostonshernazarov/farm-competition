package v1

import (
	"musobaqa/farm-competition/api/models"
	"musobaqa/farm-competition/internal/entity"
	l "musobaqa/farm-competition/internal/pkg/logger"
	"musobaqa/farm-competition/internal/pkg/otlp"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

// CREATE ANIMAL GIVEN EATABLES
// @Summary CREATE ANIMAL GIVEN EATABLES
// @Description Api for Create animal given eatables
// @Tags GIVEN-EATABLES
// @Accept json
// @Produce json
// @Param Given-Eatables body models.AnimaGivenEatablesReq true "createModel"
// @Success 201 {object} models.AnimaGivenEatablesRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/given-eatables [post]
func (h *HandlerV1) CreateGivenEatables(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimaGivenEatablesReq
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

	eatablesRes, err := h.Feeding.Create(ctx, &entity.Feeding{
		AnimalID:   body.AnimalID,
		EatablesID: body.EatablesID,
		Category:   body.Category,
		Daily:      dailyReq,
		Day:        body.Day,
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

	c.JSON(http.StatusCreated, &models.AnimaGivenEatablesRes{
		ID:         eatablesRes.ID,
		AnimalID:   eatablesRes.AnimalID,
		EatablesID: eatablesRes.Eatables.ID,
		Daily:      res,
		Category:   eatablesRes.Category,
		Day:        eatablesRes.Day,
	})
}

// UPDATE ANIMAL GIVEN EATABLES
// @Summary UPDATE ANIMAL GIVEN EATABLES
// @Description Api for Update animal given eatables
// @Tags GIVEN-EATABLES
// @Accept json
// @Produce json
// @Param Given-Eatables body models.AnimaGivenEatablesRes true "UpdateModel"
// @Success 200 {object} models.AnimaGivenEatablesRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/given-eatables [put]
func (h *HandlerV1) UpdateGivenEatables(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateAnimalEatablesInfo")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimaGivenEatablesRes
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

	res, err := h.Feeding.Update(ctx, &entity.Feeding{
		ID:         body.ID,
		AnimalID:   body.AnimalID,
		EatablesID: body.EatablesID,
		Category:   body.Category,
		Daily:      dailyReq,
		Day:        body.Day,
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

	c.JSON(http.StatusOK, &models.AnimaGivenEatablesRes{
		ID:         res.ID,
		AnimalID:   res.AnimalID,
		EatablesID: res.Eatables.ID,
		Daily:      resDaily,
		Category:   res.Category,
		Day: res.Day,
	})
}


// DELETE ANIMAL GIVEN EATABLES
// @Summary DELETE ANIMAL GIVEN EATABLES
// @Description Api for Delete animal given eatables
// @Tags GIVEN-EATABLES
// @Accept json
// @Produce json
// @Param id path string true "Animal Given Eatables ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/given-eatables/{id} [delete]
func (h *HandlerV1) DeleteGivenEatables(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	err := h.Feeding.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete animal given product", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Result{
		Message: "Animal given eatables has been deleted",
	})
}
