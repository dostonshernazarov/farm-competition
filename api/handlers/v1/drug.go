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

// CREATE DRUG
// @Summary CREATE DRUG
// @Description Api for Create new drug
// @Tags DRUG
// @Accept json
// @Produce json
// @Param Drug body models.DrugReq true "createModel"
// @Success 201 {object} models.DrugRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/drugs [post]
func (h *HandlerV1) CreateDrug(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "CreateDrug")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.DrugReq
		res *entity.Drug
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

	resName, err := h.Drug.UniqueDrugName(ctx, body.DrugName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	if resName != 0 {
		getDrug, err := h.Drug.Get(ctx, map[string]string{"name": body.DrugName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}

		res, err = h.Drug.Update(ctx, &entity.Drug{
			ID:          getDrug.ID,
			Name:        getDrug.Name,
			Status:      getDrug.Status,
			Capacity:    getDrug.Capacity + uint64(body.TotalCapacity),
			Union:       getDrug.Union,
			Description: getDrug.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Message: models.InternalMessage,
			})
			h.Logger.Error(err.Error())
			return
		}
	} else {

		res, err = h.Drug.Create(ctx, &entity.Drug{
			ID:          uuid.New().String(),
			Name:        body.DrugName,
			Status:      body.Status,
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

	c.JSON(http.StatusCreated, &models.DrugRes{
		Id:            res.ID,
		DrugName:      res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.Capacity),
		Status:        res.Status,
	})
}

// GET DRUG
// @Summary GET DRUG BY DRUG ID
// @Description Api for Get drug by ID
// @Tags DRUG
// @Accept json
// @Produce json
// @Param id path string true "Drug ID"
// @Success 200 {object} models.DrugRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/drugs/{id} [get]
func (h *HandlerV1) GetDrug(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "GetDrug")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	res, err := h.Drug.Get(ctx, map[string]string{"id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.DrugRes{
		Id:            res.ID,
		DrugName:      res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.Capacity),
		Status:        res.Status,
	})
}

// LIST DRUG
// @Summary LIST DRUG
// @Description Api for ListDrug by page limit and extra values
// @Tags DRUG
// @Accept json
// @Produce json
// @Param request query models.Pagination true "request"
// @Param request query models.DrugFieldValues true "request"
// @Success 200 {object} models.ListDrugsRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/drugs [get]
func (h *HandlerV1) ListDrug(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "ListDrug")
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
	union := c.Query("union")
	status := c.Query("status")

	mapD := map[string]interface{}{
		"name":   name,
		"union":  union,
		"status": status,
	}

	res, err := h.Drug.List(ctx, params.Page, params.Limit, mapD)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	var resList []*models.DrugRes
	for _, i := range res.Drugs {
		var resItem models.DrugRes
		resItem.Id = i.ID
		resItem.DrugName = i.Name
		resItem.Description = i.Description
		resItem.Union = i.Union
		resItem.TotalCapacity = int64(i.Capacity)
		resItem.Status = i.Status

		resList = append(resList, &resItem)
	}

	c.JSON(http.StatusOK, &models.ListDrugsRes{
		Drugs: resList,
	})
}

// UPDATE
// @Summary UPDATE DRUG
// @Description Api for Update drug by drug id
// @Tags DRUG
// @Accept json
// @Produce json
// @Param Drug body models.DrugRes true "createModel"
// @Success 200 {object} models.DrugRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/drugs [put]
func (h *HandlerV1) UpdateDrug(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "UpdateDrug")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.DrugRes
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	res, err := h.Drug.Update(ctx, &entity.Drug{
		ID:          body.Id,
		Name:        body.DrugName,
		Status:      body.Status,
		Capacity:    uint64(body.TotalCapacity),
		Union:       body.Union,
		Description: body.Description,
	})
	if err != nil {
		c.JSON(500, models.InternalMessage)
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.DrugRes{
		Id:            res.ID,
		DrugName:      res.Name,
		Union:         res.Union,
		Description:   res.Description,
		TotalCapacity: int64(res.Capacity),
		Status:        res.Status,
	})
}

// DELETE
// @Summary DELETE DRUG
// @Description Api for Delete drug by drug ID
// @Tags DRUG
// @Accept json
// @Produce json
// @Param id path string true "Drug ID"
// @Success 200 {object} models.Result
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/drugs/{id} [delete]
func (h *HandlerV1) DeleteDrug(c *gin.Context) {
	ctx, span := otlp.Start(c, "api", "DeleteDrug")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	id := c.Param("id")

	_, err := h.Drug.Get(ctx, map[string]string{"id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.NotAvailable,
		})
		h.Logger.Error("failed to get drug in delete", l.Error(err))
		return
	}

	err = h.Drug.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.InternalMessage)
		h.Logger.Error("failed to delete product", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Result{
		Message: "Drug` has been deleted",
	})
}
