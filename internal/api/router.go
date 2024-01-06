package api

import (
	"html/template"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/gvidow/go-technical-equipment/docs"

	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/redis"
	"github.com/gvidow/go-technical-equipment/internal/pkg/middlewares"
	"github.com/gvidow/go-technical-equipment/internal/pkg/service"
)

func New(cfg *config.Config, s *service.Service, tmpl *template.Template, blacklist *redis.Client) *gin.Engine {
	gin.SetMode(cfg.Mode)
	r := gin.Default()
	r.Static("/upload", "./upload")
	produceRouting(r, s, cfg, blacklist)
	return r
}

func produceRouting(r *gin.Engine, s *service.Service, cfg *config.Config, bl *redis.Client) {
	r.Use(middlewares.Auth(cfg.JWT, bl))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api/v1/")
	{
		eq := api.Group("/equipment")
		{
			eq.GET("/list/active", s.GetListEquipments)
			eq.GET("/list", s.FeedEquipment)
			eq.GET("/get/:id", s.GetOneEquipment)

			eq.Use(middlewares.RequireAuth()).POST("/last/:id", s.AddEquipmentInLastRequest)

			sub := eq.Group("", middlewares.RequireAuth(ds.Moderator))
			{
				sub.POST("/add", s.AddNewEquipment)
				sub.PUT("/edit/:id", s.EditEquipment)
				sub.DELETE("/delete/:id", s.DeleteEquipment)
			}
		}

		req := api.Group("/request")
		{
			req.Use(middlewares.RequireAuth())

			req.GET("/get/:id", s.GetRequest)
			req.GET("/list", s.ListRequest)

			req.PUT("/format/:id", s.OperationRequest)
			req.DELETE("/delete/:id", s.DropRequest)

			req.Group("/status/change/moderator/:id").
				Use(middlewares.RequireAuth(ds.Moderator)).
				PUT("/", s.StatusChangeByModerator)
		}

		order := api.Group("/order").Use(middlewares.RequireAuth())
		{
			order.PUT("/edit/count/:id", s.EditCount)
			order.DELETE("/delete/:id", s.DeleteOrder)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", s.Login)
			auth.POST("/signup", s.Signup)
			auth.Use(middlewares.RequireAuth()).DELETE("/logout", s.Logout)
		}
	}
}
