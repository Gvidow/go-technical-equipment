package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
)

func (s *Service) GetListEquipments(c *gin.Context) {
	title := c.Request.FormValue("title")
	if title == "" {
		equipments, err := s.u.GetListEquipments()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список оборудования"})
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "body": equipments})
	} else {
		equipments, err := s.u.GetListEquipmentsWithFilter(title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список оборудования"})
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "body": equipments})
	}
}

func (s *Service) GetOneEquipment(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в пути запроса отсутсвует id оборудования"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "id оборудования в пути запроса должен быть натуральным числом"})
		return
	}

	equipment, err := s.u.GetOneEquipmentByID(int(id))
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

	err = s.u.AddNewEquipment(c.Request.Context(), title, description, f, fh.Header.Get("Content-Type"), fh.Size, fh.Filename)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "bad request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func (s *Service) EditEquipment(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в пути запроса отсутсвует id оборудования"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "id оборудования в пути запроса должен быть натуральным числом"})
		return
	}

	equipment, err := s.u.GetOneEquipmentByID(int(id))
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
		fileURL, err := s.u.PutFileInMinio(c.Request.Context(), file, fh.Header.Get("Content-Type"), fh.Size, fh.Filename)
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

	err = s.u.EditEquipment(equipment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось изменить оборудование"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (s *Service) DeleteEquipment(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "в пути запроса отсутсвует id оборудования"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "id оборудования в пути запроса должен быть натуральным числом"})
		return
	}

	err = s.u.DeleteEquipmentByID(int(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "оборудование не нашлось"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (s *Service) AddEquipmentInLastRequest(c *gin.Context) {
	r := c.Request.Context().Value(middlewares.ContextUserID)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": fmt.Sprintf("call method from user with id=%v the add equipment", r)})
}
