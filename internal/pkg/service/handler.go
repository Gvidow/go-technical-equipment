package service

import (
	"fmt"
	"net/http"
	"net/url"
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
			s.log.Error(err)
		}
	} else {
		equipments, err = s.u.Equipment().GetAllEquipments()
		if err != nil {
			s.log.Error(err)
		}
	}
	data["Equipments"] = equipments

	if req, err := s.u.Request().GetEnteredRequest(); err != nil {
		s.log.Error(err)
		data["CartID"] = 0
	} else {
		data["CartID"] = req.ID
	}

	c.HTML(http.StatusOK, "index.html", data)
}

func (s *Service) Equipment(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		s.log.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	equipment, err := s.u.Equipment().GetByID(int(id))
	if err != nil {
		s.log.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "equipment.html", equipment)
}

func (s *Service) DeleteEquipment(c *gin.Context) {
	c.Request.ParseForm()
	paramID := c.Param("id")
	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		s.log.Warn("handler DeleteEquipment couldn't parse the id parameter")
	} else if s.u.Equipment().DeleteEquipmentByID(int(id)) != nil {
		s.log.Error(err)
	}
	if title, ok := c.Request.Form["title"]; ok {
		c.Redirect(http.StatusFound,
			fmt.Sprintf("%s?title=%s", MainPageURL, url.PathEscape(title[0])))
	} else {
		c.Redirect(http.StatusFound, MainPageURL)
	}
}

func (s *Service) RequestDetail(c *gin.Context) {
	idString := c.Param("id")
	requestID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		s.log.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	request, err := s.u.Request().GetRequestByID(int(requestID))
	if err != nil {
		s.log.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// if q, ok := c.Request.Form["title"]; ok {
	// 	data["Search"] = q[0]
	// 	equipments, err = s.u.Equipment().SearchEquipmentsByTitle(q[0])
	// 	if err != nil {
	// 		s.log.Error(err)
	// 	}
	// } else {
	// 	equipments, err = s.u.Equipment().GetAllEquipments()
	// 	if err != nil {
	// 		s.log.Error(err)
	// 	}
	// }
	// data["Equipments"] = equipments

	// if req, err := s.u.Request().GetEnteredRequest(); err != nil {
	// 	s.log.Error(err)
	// 	data["CartID"] = 0
	// } else {
	// 	data["CartID"] = req.ID
	// }
	c.HTML(http.StatusOK, "request.html", request)
}

func (s *Service) BadRequest(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, MainPageURL)
}
