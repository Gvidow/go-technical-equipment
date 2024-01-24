package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/request"
	mw "github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

// ShowAccount godoc
// @Summary      Show an requests
// @Description  get requests of users
// @Tags         request
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        creatorProfile       query      string     false  "Creator username filter" example(user)
// @Param        moderatorProfile     query      string     false  "Moderator username filter" example(admin)
// @Param        status               query      string     false  "Status filter" example(completed)
// @Param        formatedAfter        query      string     false  "Request formatted after filter" format(date) example(30.12.2023)
// @Param        formatedBefore       query      string     false  "Request formatted before filter" format(date) example(30.12.2023)
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /request/list [get]
func (s *Service) ListRequest(c *gin.Context) {
	ctxUser := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)

	user := &ds.User{ID: ctxUser.UserID}
	user.SetRole(ctxUser.Role)

	feedCfg, err := encodeFeedRequestConfig(c.Request.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ошибка в параметрах запроса"})
		return
	}

	feed, err := s.reqCase.GetFeedRequests(feedCfg, user)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список заявок"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "body": feed})
}

// ShowAccount godoc
// @Summary      Show an request in detail
// @Description  get information about the request
// @Tags         request
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id               path      int     true  "Request id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /request/get/{id} [get]
func (s *Service) GetRequest(c *gin.Context) {
	ctxUser := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)

	user := &ds.User{ID: ctxUser.UserID}
	user.SetRole(ctxUser.Role)

	requestID, err := FetchIdFromURLPath(c)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось извлечь id заявки из запроса"})
		return
	}

	request, err := s.reqCase.GetRequestByID(requestID, user)
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

// ShowAccount godoc
// @Summary      Formation of the request
// @Description  forming an request in the status 'entered' with equipment
// @Tags         request
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id               path      int     true  "Request id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /request/format/{id} [put]
func (s *Service) OperationRequest(c *gin.Context) {
	userID := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole).UserID
	requestID, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "для формирования заявки требуется передать в пути её id"})
		return
	}

	err = s.reqCase.ToFormRequest(requestID, userID)
	if err == request.ErrStatusCannotChange {
		s.log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "статус заявки не позволяет совершать её формирование"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось сформировать заявку"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "заявка успешно сформирована"})
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

// ShowAccount godoc
// @Summary      Change status of the request
// @Description  edit status request from 'entered' to received
// @Tags         request
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id               path      int     true  "Request id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /request/status/change/moderator/{id} [put]
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

	err = s.reqCase.StatusChangeByModerator(user.UserID, requestID, newStatus.Status)

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
		s.log.Info(err.Error())
		message = "не удалось изменить статус"
		status = http.StatusNotFound
	}
	c.JSON(status, gin.H{"status": "error", "message": message})
}

// ShowAccount godoc
// @Summary      Delete of the request
// @Description  delete an request in the status 'entered' with equipment
// @Tags         request
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id               path      int     true  "Request id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /request/delete/{id} [delete]
func (s *Service) DropRequest(c *gin.Context) {
	ctxUser := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)

	user := &ds.User{ID: ctxUser.UserID}
	user.SetRole(ctxUser.Role)

	id, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	err = s.reqCase.DropRequest(id, user)
	if err == request.ErrStatusCannotChange {
		s.log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "статус заявки не позволяет совершать её удаление"})
		return
	}
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось удалить заявку"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}

}
