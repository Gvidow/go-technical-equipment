package app

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/api"
	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/equipment"
	"github.com/gvidow/go-technical-equipment/internal/app/usecase"
	"github.com/gvidow/go-technical-equipment/internal/pkg/service"
	"github.com/gvidow/go-technical-equipment/logger"
)

type Application struct {
	log    *logger.Logger
	cfg    *config.Config
	router *gin.Engine
}

func New(log *logger.Logger, cfg *config.Config) (*Application, error) {
	// db, err := gorm.Open(postgres.Open(dsn.FromEnv()))
	// if err != nil {
	// 	return nil, err
	// }
	// repo := equipment.NewRepository(db)
	repo := equipment.NewStorageRepository()
	u := usecase.New(repo)
	tmpl := template.Must(template.ParseGlob("templates/*"))
	s := service.New(log, tmpl, u)
	r := api.New(cfg, s, tmpl)

	return &Application{
		log:    log,
		cfg:    cfg,
		router: r,
	}, nil
}

func (a *Application) Run() error {
	a.log.Info(fmt.Sprintf("start server on %s:%s with mode=%s", a.cfg.ServiceHost, a.cfg.ServicePort, a.cfg.Mode))
	return a.router.Run(a.cfg.ServiceHost + ":" + a.cfg.ServicePort)
}
