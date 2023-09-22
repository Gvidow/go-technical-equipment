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
	r.GET("/equipment/:id", s.Equipment)
	r.NoRoute(s.BadRequest)
}
