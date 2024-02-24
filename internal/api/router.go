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
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "./static")
	r.Static("/upload", "./upload")
	produceRouting(r, s)
	return r
}

func produceRouting(r *gin.Engine, s *service.Service) {
	r.GET(service.MainPageURL, s.MainPage)
	eq := r.Group("/equipment")
	{
		eq.GET("/:id", s.Equipment)
		eq.POST("/:id", s.DeleteEquipment)
	}

	req := r.Group("/request")
	{
		req.GET("/:id", s.RequestDetail)
	}

	r.NoRoute(s.BadRequest)
}
