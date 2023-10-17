package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) GetListEquipments(c *gin.Context) {
	equipments, err := s.u.GetListEquipments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось получить список оборудования"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "body": equipments})
}

func (s *Service) GetOneEquipment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"handler": "GetOneEquipment", "path": c.Request.URL.Path, "method": c.Request.Method})
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
	c.JSON(http.StatusOK, gin.H{"handler": "EditEquipment", "path": c.Request.URL.Path, "method": c.Request.Method})
}

func (s *Service) DeleteEquipment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"handler": "DeleteEquipment", "path": c.Request.URL.Path, "method": c.Request.Method})
}

func (s *Service) AddEquipmentInLastRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"handler": "AddEquipmentInLastRequest", "path": c.Request.URL.Path, "method": c.Request.Method})
}
