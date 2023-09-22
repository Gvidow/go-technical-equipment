package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
)

var MainPageURL = "/main"

func (s *Service) MainPage(c *gin.Context) {
	c.Request.ParseForm()
	var equipments []ds.Equipment
	var err error
	data := gin.H{}
	if q, ok := c.Request.Form["title"]; ok {
		data["Search"] = q[0]
		equipments, err = s.u.Equipment().SearchEquipmentsByTitle(q[0])
		if err != nil {
			s.log.Error(err.Error())
		}
	} else {
		equipments, err = s.u.Equipment().GetAllEquipments()
		if err != nil {
			s.log.Error(err.Error())
		}
	}
	data["Equipments"] = equipments
	c.HTML(http.StatusOK, "index.html", data)
}

func (s *Service) Equipment(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		s.log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	equipment, err := s.u.Equipment().GetByID(int(id))
	if err != nil {
		s.log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "equipment.html", equipment)
}

func (s *Service) BadRequest(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, MainPageURL)
}
