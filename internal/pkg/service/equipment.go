package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	mw "github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

func (s *Service) GetListEquipments(c *gin.Context) {
	r := c.Request.Context().Value(mw.ContextUserID)
	var (
		lastRequest *ds.Request
		err         error
	)

	if userID, ok := r.(int); ok {
		lastRequest, err = s.reqCase.GettingUserLastRequest(userID)
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

func (s *Service) FeedEquipment(c *gin.Context) {
	userID, ok := c.Request.Context().Value(mw.ContextUserID).(int)
	var (
		lastRequest *ds.Request
		err         error
	)

	if ok {
		lastRequest, err = s.reqCase.GettingUserLastRequest(userID)
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

func (s *Service) AddEquipmentInLastRequest(c *gin.Context) {
	userID, ok := c.Request.Context().Value(mw.ContextUserID).(int)
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

	req, err := s.reqCase.GettingUserLastRequest(userID)
	if err != nil {
		req, err = s.reqCase.CreateDraftRequest(userID)
		if err != nil {
			s.log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось создать заявку"})
			return
		}
	}

	s.log.Sugar().Infof("add equipment(%d) in request(%d)", equipmentID, req.Id())

	err = s.orCase.AddEquipmentInRequest(equipmentID, req.Id())
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось добавить оборудование в заявку"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
