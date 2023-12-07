package service

import (
	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/auth"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/equipment"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/order"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/request"
	"github.com/gvidow/go-technical-equipment/logger"
)

type Service struct {
	log      *logger.Logger
	cfg      *config.Config
	eqCase   *equipment.Usecase
	reqCase  *request.Usecase
	orCase   *order.Usecase
	authCase *auth.Usecase
}

func New(log *logger.Logger, cfg *config.Config, eqCase *equipment.Usecase, reqCase *request.Usecase, orCase *order.Usecase, authCase *auth.Usecase) *Service {
	return &Service{log, cfg, eqCase, reqCase, orCase, authCase}
}
