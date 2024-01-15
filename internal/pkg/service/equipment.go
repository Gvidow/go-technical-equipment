package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/order"
	mw "github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

func (s *Service) GetListEquipments(c *gin.Context) {
	r := c.Request.Context().Value(mw.ContextUser)
	var (
		lastRequest *ds.Request
		err         error
	)

	if user, ok := r.(mw.UserWithRole); ok {
		lastRequest, err = s.reqCase.GettingUserLastRequest(user.UserID)
	}

	if err != nil {
		s.log.Error(err)
	}

	title := c.Request.FormValue("title")
	if title == "" {
		equipments, err := s.eqCase.GetListEquipments()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список оборудования"})
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "body": gin.H{"equipments": equipments, "last_request_id": lastRequest.Id()}})
	} else {
		equipments, err := s.eqCase.GetListEquipmentsWithFilter(title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список оборудования"})
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "body": gin.H{"equipments": equipments, "last_request_id": lastRequest.Id()}})
	}
}

// ShowAccount godoc
// @Summary      Show an equipments
// @Description  get list equipments
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Param        equipment      query      string  false  "Title filter"
// @Param        status         query      string  false  "Status filter"
// @Param        createdAfter   query      string  false  "Created after" format(date) example(30.12.2023)
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /equipment/list [get]
func (s *Service) FeedEquipment(c *gin.Context) {
	user, ok := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)
	var (
		lastRequest *ds.Request
		err         error
	)

	if ok {
		lastRequest, err = s.reqCase.GettingUserLastRequest(user.UserID)
	}

	if err != nil {
		s.log.Error(err)
	}

	cfg, err := encodeFeedConfig(c.Request.URL)
	if err != nil {
		s.log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ошибка в указании параметров запроса"})
		return
	}
	equipments, err := s.eqCase.ViewFeedEquipment(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список оборудования"})
		s.log.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "body": gin.H{"equipments": equipments, "last_request_id": lastRequest.Id()}})
}

// ShowAccount godoc
// @Summary      Show an equipment in detail
// @Description  get information about the equipment
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Param        id         path      int  true  "Equipment id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /equipment/get/{id} [get]
func (s *Service) GetOneEquipment(c *gin.Context) {
	id, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	equipment, err := s.eqCase.GetOneEquipmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "оборудование не нашлось"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "body": equipment})
}

// ShowAccount godoc
// @Summary      Adding a equipment
// @Description  adding new equipment to the turnover
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id         path      int  true  "Equipment id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /equipment/add [post]
func (s *Service) AddNewEquipment(c *gin.Context) {
	err := c.Request.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		s.log.Info("bad request")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "bad request"})
		return
	}
	defer c.Request.Body.Close()

	title := c.Request.FormValue("title")
	description := c.Request.FormValue("description")

	f, fh, err := c.Request.FormFile("picture")
	if err != nil {
		s.log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "bad request"})
		return
	}
	defer f.Close()

	newID, err := s.eqCase.AddNewEquipment(c.Request.Context(), title, description, f, fh.Header.Get("Content-Type"), fh.Size, fh.Filename)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "bad request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "body": map[string]int{"id": newID}})
}

// ShowAccount godoc
// @Summary      Update a equipment
// @Description  edit the active equipment
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id         path      int  true  "Equipment id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /equipment/edit/{id} [put]
func (s *Service) EditEquipment(c *gin.Context) {
	id, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	equipment, err := s.eqCase.GetOneEquipmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "оборудование не нашлось"})
		return
	}

	err = c.Request.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось распасить тело запроса"})
		return
	}
	defer c.Request.Body.Close()

	count := 0
	if countStr := c.Request.FormValue("count"); countStr != "" {
		countInt64, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil || countInt64 <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "count должен быть натуральным числом"})
			return
		}
		count = int(countInt64)
	}

	file, fh, err := c.Request.FormFile("picture")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удаётся прочитать файл"})
	} else if err != http.ErrMissingFile {
		defer file.Close()
		fileURL, err := s.eqCase.PutFileInMinio(c.Request.Context(), file, fh.Header.Get("Content-Type"), fh.Size, fh.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		}
		equipment.Picture = fileURL
	}

	if title := c.Request.FormValue("title"); title != "" {
		equipment.Title = title
	}

	if description := c.Request.FormValue("description"); description != "" {
		equipment.Description = description
	}

	if count != 0 {
		equipment.Count = count
	}

	err = s.eqCase.EditEquipment(equipment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось изменить оборудование"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// ShowAccount godoc
// @Summary      Delete a equipment
// @Description  change the equipment status from 'active' to 'delete'
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id         path      int  true  "Equipment id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /equipment/delete/{id} [delete]
func (s *Service) DeleteEquipment(c *gin.Context) {
	id, err := FetchIdFromURLPath(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	err = s.eqCase.DeleteEquipmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "оборудование не нашлось"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// ShowAccount godoc
// @Summary      Adding equipment to the shopping cart
// @Description  add a service to the user's draft request
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Param        id         path      int  true  "Equipment id"
// @Success      200  {object}  int
// @Failure      400  {object}  int
// @Failure      404  {object}  string
// @Failure      500  {object}  int
// @Router       /equipment/last/{id} [post]
func (s *Service) AddEquipmentInLastRequest(c *gin.Context) {
	user, ok := c.Request.Context().Value(mw.ContextUser).(mw.UserWithRole)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "для добавления оборудования в корзину нужно авторизоваться"})
		return
	}

	equipmentID, err := FetchIdFromURLPath(c)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "в пути запроса должен быть указан id оборудования - натуральное число"})
		return
	}

	req, err := s.reqCase.GettingUserLastRequest(user.UserID)
	if err != nil {
		req, err = s.reqCase.CreateDraftRequest(user.UserID)
		if err != nil {
			s.log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось создать заявку"})
			return
		}
	}

	s.log.Sugar().Infof("add equipment(%d) in request(%d)", equipmentID, req.Id())

	err = s.orCase.AddEquipmentInRequest(equipmentID, req.Id())
	if err == order.ErrEquipmentNotFound {
		s.log.Info("equipment is deleted")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "заявка не найдена"})
		return
	}

	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось добавить оборудование в заявку"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
