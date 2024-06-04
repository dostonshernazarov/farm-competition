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

// CREATE ANIMAL PRODUCT
// @Summary CREATE ANIMAL PRODUCT
// @Description Api for Create product which has got from animal
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param Animal-Product body models.DrugReq true "createModel"
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

	res, err := h.Drug.Create(ctx, &entity.Drug{

	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: models.InternalMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, &models.AnimalProductRes{
		Id:        res.ID,
		AnimalID:  "",
		ProductName: "",
		Capacity:  0,
		GetTime:   "",
	})
}

// GET DRUG
// @Summary GET ANIMAL PRODUCT BY ID
// @Description Api for Get product which has got from animal by ID
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param id path string true "Animal-Product ID"
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

	res, err := h.Drug.Get(ctx, map[string]string{"id":id})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.AnimalProductRes{
		Id:          id,
		AnimalID:    res.ID,
		ProductName: "",
		Capacity:    0,
		GetTime:     "",
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
	// ctx, span := otlp.Start(c, "api", "ListAnimalProducts")
	// span.SetAttributes(
	// 	attribute.Key("method").String(c.Request.Method),
	// 	attribute.Key("host").String(c.Request.Host),
	// )
	// defer span.End()

	// queryParams := c.Request.URL.Query()
	// params, errStr := utils.ParseQueryParam(queryParams)
	// if errStr != nil {
	// 	c.JSON(http.StatusBadRequest, models.Error{
	// 		Message: models.WrongInfoMessage,
	// 	})
	// 	return
	// }


	// animalID := c.Query("animal_id")

	// mapA := map[string]interface{}{
	// 	"animal_id":  animalID,
	// }

	// res, err := h.Drug.List(ctx, params.Page, params.Limit, mapA)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, models.Error{
	// 		Message: models.InternalMessage,
	// 	})
	// 	h.Logger.Error(err.Error())
	// 	return
	// }

	// var resList []*models.AnimalProductRes
	// for _, i := range res.AnimalProducts {
	// 	var resItem models.AnimalProductRes
	// 	resItem.Id = i.ID
	// 	resItem.AnimalProductsName = i.Name
	// 	resItem.Description = i.Description
	// 	resItem.Union = i.Union
	// 	resItem.TotalCapacity = int64(i.Capacity)
	// 	resItem.Status = i.Status

	// 	resList = append(resList, &resItem)
	// }

	c.JSON(http.StatusOK, &models.ListAnimalProductsRes{
		AnimalProducts: []*models.AnimalProductRes{},
	})
}


// UPDATE
// @Summary UPDATE ANIMAL PRODUCT
// @Description Api for Update product which has got from animal
// @Tags ANIMAL-PRODUCT
// @Accept json
// @Produce json
// @Param Animal-Product body models.AnimalProductRes true "createModel"
// @Success 200 {object} models.AnimalProductRes
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/animals/products [put]
func (h *HandlerV1) UpdateAnimalProduct(c *gin.Context) {
	_, span := otlp.Start(c, "api", "UpdateAnimalProduct")
	span.SetAttributes(
		attribute.Key("method").String(c.Request.Method),
		attribute.Key("host").String(c.Request.Host),
	)
	defer span.End()

	var (
		body models.AnimalProductRes
	)


	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: models.WrongInfoMessage,
		})
		h.Logger.Error("failed to bind json", l.Error(err))
		return
	}

	// res, err := h.Drug.Update(ctx, &entity.Drug{
	// 	ID:          body.Id,
	// 	Name:        body.DrugName,
	// 	Status:      body.Status,
	// 	Capacity:    uint64(body.TotalCapacity),
	// 	Union:       body.Union,
	// 	Description: body.Description,
	// })
	// if err != nil {
	// 	c.JSON(500, models.InternalMessage)
	// 	h.Logger.Error(err.Error())
	// 	return
	// }

	c.JSON(http.StatusOK, &models.AnimalProductRes{

	})
}

// DELETE
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

	_, err := h.Drug.Get(ctx, map[string]string{})
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
		Message: "Product from animal has been deleted",
	})
}