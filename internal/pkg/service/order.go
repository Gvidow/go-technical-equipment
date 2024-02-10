package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	mw "github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

type bodyCount struct {
	Count int
}

// ShowAccount godoc
// @Summary      Edit count equipment from the request
// @Description  edit count equipment from the user's request with status 'entered'
// @Tags         request
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Equipment ID"
// @Param        new_count   body      bodyCount  true  "New count for equipment with id"
// @Success      200  {object}  ResponseOk
// @Failure      400  {object}  ResponseError
// @Failure      404  {object}  ResponseError
// @Failure      500  {object}  ResponseError
// @Router       /order/edit/count/{id} [put]
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
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в пути должен быть указан id оборудования - натуральное число"})
		return
	}

	var count bodyCount
	err = json.NewDecoder(c.Request.Body).Decode(&count)
	c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в запросе должно быть передано новое количество - натуральное число"})
		return
	}

	countNew := count.Count

	if countNew < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "количество должно быть положительным числом"})
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

// ShowAccount godoc
// @Summary      Removing equipment from the request
// @Description  delete equipment from the user's request with status 'entered'
// @Tags         request
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id               path      int     true  "Equipment id"
// @Success      200  {object}  ResponseOk
// @Failure      400  {object}  ResponseError
// @Failure      404  {object}  ResponseError
// @Failure      500  {object}  ResponseError
// @Router       /order/delete/{id} [delete]
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
