package service

import (
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/equipment"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/order"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/request"
	"github.com/gvidow/go-technical-equipment/logger"
)

type Service struct {
	log     *logger.Logger
	eqCase  *equipment.Usecase
	reqCase *request.Usecase
	orCase  *order.Usecase
}

func New(log *logger.Logger, eqCase *equipment.Usecase, reqCase *request.Usecase, orCase *order.Usecase) *Service {
	return &Service{log, eqCase, reqCase, orCase}
}
