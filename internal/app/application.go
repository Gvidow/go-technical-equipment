package app

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gvidow/go-technical-equipment/internal/api"
	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/app/dsn"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/equipment"
	orRepo "github.com/gvidow/go-technical-equipment/internal/app/repository/order"
	reqRepo "github.com/gvidow/go-technical-equipment/internal/app/repository/request"
	userRepo "github.com/gvidow/go-technical-equipment/internal/app/repository/user"
	ucEquipment "github.com/gvidow/go-technical-equipment/internal/app/usecases/equipment"
	orCase "github.com/gvidow/go-technical-equipment/internal/app/usecases/order"
	reqCase "github.com/gvidow/go-technical-equipment/internal/app/usecases/request"
	"github.com/gvidow/go-technical-equipment/internal/pkg/service"
	"github.com/gvidow/go-technical-equipment/logger"
)

type Application struct {
	log    *logger.Logger
	cfg    *config.Config
	router *gin.Engine
}

func New(log *logger.Logger, cfg *config.Config) (*Application, error) {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()))
	if err != nil {
		return nil, err
	}
	repo := equipment.NewRepository(db)
	u, err := ucEquipment.New(repo, ucEquipment.NewMinioConfig("http://localhost:9000",
		"minio", "minio124").SetBucket("equipment"))
	if err != nil {
		return nil, fmt.Errorf("new equipment usecase: %w", err)
	}

	tmpl := template.Must(template.ParseGlob("templates/*"))
	s := service.New(log, u, reqCase.NewUsecase(reqRepo.NewRepository(db), userRepo.NewUserRepo(db)), orCase.NewUsecase(orRepo.NewRepository(db), repo))
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
