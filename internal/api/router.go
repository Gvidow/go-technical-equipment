package api

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
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
	r.Use(middlewares.Auth())
	api := r.Group("/api/v1/")
	{
		eq := api.Group("/equipment")
		{
			eq.GET("/list/active", s.GetListEquipments)
			eq.GET("/list", s.FeedEquipment)
			eq.GET("/get/:id", s.GetOneEquipment)
			eq.POST("/add", s.AddNewEquipment)
			eq.PUT("/edit/:id", s.EditEquipment)
			eq.DELETE("/delete/:id", s.DeleteEquipment)
			eq.POST("/last/:id", s.AddEquipmentInLastRequest)
		}

		req := api.Group("/request")
		{
			req.GET("/list", s.ListRequest)
			req.GET("/get/:id", s.GetRequest)

			sub := req.Group("", middlewares.RequireAuth())
			{
				sub.PUT("/format/:id", s.OperationRequest)
				sub.PUT("/status/change/moderator/:id", s.StatusChangeByModerator)
				sub.DELETE("/delete/:id", s.DropRequest)
			}
		}

		order := api.Group("/order")
		{
			order.Use(middlewares.RequireAuth())
			order.PUT("/edit/count/:id", s.EditCount)
			order.DELETE("/delete/:id", s.DeleteOrder)
		}
	}
	// r.NoRoute(s.BadRequest)
}
