package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	mw "github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /edit/{id} [get]
func (s *Service) EditCount(c *gin.Context) {
	user := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)

	request, err := s.reqCase.GettingUserLastRequest(user.UserID)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "возникла проблема с получение черновой заявки пользователя"})
		return
	}

	equipmentID, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в пути должно быть передано новое количество - натуральное число"})
		return
	}

	countNew, err := strconv.ParseInt(c.Request.FormValue("count"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в пути должен быть указан id оборудования - натуральное число"})
		return
	}

	err = s.orCase.EditCountEquipmentsInRequest(equipmentID, request.ID, int(countNew))
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "не нашлось такого оборудования в заданной заявке"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "колличество оборудования в заявке успешно изменено"})
}

func (s *Service) DeleteOrder(c *gin.Context) {
	user := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)

	request, err := s.reqCase.GettingUserLastRequest(user.UserID)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "возникла проблема с получение черновой заявки пользователя"})
		return
	}

	equipmentID, err := FetchIdFromURLPath(c)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	err = s.orCase.DeleteEquipmentFromRequest(equipmentID, request.ID)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ошибка при удалении оборудования из заявки"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "оборудование успешно удаленно из заявки"})
}
