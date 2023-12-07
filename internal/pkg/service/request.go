package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/request"
	mw "github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

func (s *Service) ListRequest(c *gin.Context) {
	feedCfg, err := encodeFeedRequestConfig(c.Request.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ошибка в параметрах запроса"})
		return
	}

	feed, err := s.reqCase.GetFeedRequests(feedCfg)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список заявок"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "body": feed})
}

func (s *Service) GetRequest(c *gin.Context) {
	requestID, err := FetchIdFromURLPath(c)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось извлечь id заявки из запроса"})
		return
	}

	request, err := s.reqCase.GetRequestByID(requestID)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "не удалось найти выбранную заявку"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "заявка успешно найдена", "body": request})
}

func (s *Service) EditRequest(c *gin.Context) {
	requestID, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "для изменения заявки требуется передать в пути её id"})
		return
	}

	requestUpdate := make(map[string]any)
	err = json.NewDecoder(c.Request.Body).Decode(&requestUpdate)
	defer c.Request.Body.Close()
	if err != nil {
		s.log.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось распарсить тело запроса"})
	}

	err = s.reqCase.EditRequest(requestID, requestUpdate)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось изменить заявку"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "заявка успешно изменена"})
}

func (s *Service) StatusChangeByCreator(c *gin.Context) {
	user, ok := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "для изменения статуса заявки вы должны авторизоваться по модератором"})
		return
	}

	requestID, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "для изменения статуса заявки требуется передать в пути её id"})
		return
	}

	var newStatus struct {
		Status string
	}

	err = json.NewDecoder(c.Request.Body).Decode(&newStatus)
	defer c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ну удалось прочитать тело запроса"})
		return
	}

	err = s.reqCase.ChangeStatusRequest(user.UserID, requestID, newStatus.Status, ds.RegularUser)
	var message string
	var status int
	switch err {
	case nil:
		c.JSON(status, gin.H{"status": "ok", "message": "статус успешно изменен"})
		return
	case request.ErrNotAccess:
		message = "у пользователя нет доступа"
		status = http.StatusForbidden
	case request.ErrRoleHaveNotAccess:
		message = "пользователь с этой ролью не может изменить статус заявки на указанный"
		status = http.StatusForbidden
	default:
		message = "не удалось изменить статус"
		status = http.StatusNotFound
	}
	c.JSON(status, gin.H{"status": "error", "message": message})
}

func (s *Service) StatusChangeByModerator(c *gin.Context) {
	user, ok := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "для изменения статуса заявки вы должны авторизоваться по модератором"})
		return
	}

	requestID, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "для изменения статуса заявки требуется передать в пути её id"})
		return
	}

	var newStatus struct {
		Status string
	}

	err = json.NewDecoder(c.Request.Body).Decode(&newStatus)
	defer c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ну удалось прочитать тело запроса"})
		return
	}

	err = s.reqCase.ChangeStatusRequest(user.UserID, requestID, newStatus.Status, ds.Moderator)
	var message string
	var status int
	switch err {
	case nil:
		c.JSON(status, gin.H{"status": "ok", "message": "статус успешно изменен"})
		return
	case request.ErrNotAccess:
		message = "у пользователя нет доступа"
		status = http.StatusForbidden
	case request.ErrRoleHaveNotAccess:
		message = "пользователь с этой ролью не может изменить статус заявки на указанный"
		status = http.StatusForbidden
	default:
		message = "не удалось изменить статус"
		status = http.StatusNotFound
	}
	c.JSON(status, gin.H{"status": "error", "message": message})
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
