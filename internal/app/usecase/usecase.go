package usecase

import (
	"github.com/gvidow/go-technical-equipment/internal/app/repository/equipment"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/request"
)

type Usercase struct {
	equipmentRepo equipment.Repository
	requestRepo   request.Repository
}

func New(er equipment.Repository, req request.Repository) *Usercase {
	return &Usercase{er, req}
}

func (u *Usercase) Equipment() equipment.Repository {
	return u.equipmentRepo
}

func (u *Usercase) Request() request.Repository {
	return u.requestRepo
}
