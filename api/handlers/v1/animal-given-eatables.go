package v1

import (
	"musobaqa/farm-competition/api/models"
	"musobaqa/farm-competition/internal/entity"
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
// @Success 201 {object} models.AnimaEatablesInfoRes
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
