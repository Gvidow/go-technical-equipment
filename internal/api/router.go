package api

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/pkg/service"
)

func New(cfg *config.Config, s *service.Service, tmpl *template.Template) *gin.Engine {
	gin.SetMode(cfg.Mode)
	r := gin.Default()
	r.Static("/upload", "./upload")
	produceRouting(r, s)
	return r
}

func produceRouting(r *gin.Engine, s *service.Service) {
	api := r.Group("/api/v1/")
	{
		eq := api.Group("/equipment")
		{
			eq.GET("/list", s.GetListEquipments)
			eq.GET("/get/:id", s.GetOneEquipment)
			eq.POST("/add", s.AddNewEquipment)
			eq.PUT("/edit/:id", s.EditEquipment)
			eq.DELETE("/delete/:id", s.DeleteEquipment)
			eq.POST("/last/:id", s.AddEquipmentInLastRequest)
		}
	}
	// r.NoRoute(s.BadRequest)
}
