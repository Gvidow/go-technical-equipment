package service

import (
	"html/template"

	"github.com/gvidow/go-technical-equipment/internal/app/usecases/equipment"
	"github.com/gvidow/go-technical-equipment/logger"
)

type Service struct {
	log  *logger.Logger
	tmpl *template.Template
	u    *equipment.Usecase
}

func New(log *logger.Logger, tmpl *template.Template, u *equipment.Usecase) *Service {
	return &Service{log, tmpl, u}
}
