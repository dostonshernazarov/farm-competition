package v1

import (
	"github.com/gin-gonic/gin"
)

// CREATE
// @Summary CREATE
// @Security BearerAuth
// @Description Api for Create
// @Tags USER
// @Accept json
// @Produce json
// @Param User body models.UserReq true "createModel"
// @Success 200 {object} models.UserRes
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
// @Router /v1/users [post]
func (h *HandlerV1) Create(c *gin.Context) {
}
