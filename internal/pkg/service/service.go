package service

import (
	"html/template"

	"github.com/gvidow/go-technical-equipment/internal/app/usecase"
	"github.com/gvidow/go-technical-equipment/logger"
)

type Service struct {
	log  *logger.Logger
	tmpl *template.Template
	u    *usecase.Usercase
}

func New(log *logger.Logger, tmpl *template.Template, u *usecase.Usercase) *Service {
	return &Service{log, tmpl, u}
}
