package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

//		{
//			"status": "ok",
//			"body": {
//				"requests": [
//					{
//						"status"
//						"equipment"
//	                 "creator"
//	                 "moderator"
//						"created_at"
//						"creator_at"
//					}
//				]
//			}
//		}
func (s *Service) ListRequest(c *gin.Context) {

}

func (s *Service) ReceivingRequest(c *gin.Context) {

}

func (s *Service) EditRequest(c *gin.Context) {

}

func (s *Service) StatusChangeByCreator(c *gin.Context) {

}

func (s *Service) StatusChangeByModerator(c *gin.Context) {
	r := c.Request.Context().Value(middlewares.ContextUserID)

	if r == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "для изменения статуса заявки вы должны авторизоваться по модератором"})
		return
	}

	id, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	var m = make(map[string]string, 1)
	err = json.NewDecoder(c.Request.Body).Decode(&m)
	defer c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "плохо"})
	}
	if newStatus, ok := m["status"]; ok {
		if newStatus != "completed" && newStatus != "canceled" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "неверный статус"})
			return
		}
		err = s.reqCase.ChangeStatusRequest(id, newStatus, "operation")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "неудалось изменить статус"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в теле запроса не указан, либо указан неправильно новый статус заявки"})
	}
}

func (s *Service) DropRequest(c *gin.Context) {
	id, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	err = s.reqCase.DropRequest(id)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось удалить заявку"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}

}
