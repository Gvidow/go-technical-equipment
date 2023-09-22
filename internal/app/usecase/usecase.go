package usecase

import (
	"github.com/gvidow/go-technical-equipment/internal/app/repository/equipment"
)

type Usercase struct {
	equipmentRepo equipment.Repository
}

func New(er equipment.Repository) *Usercase {
	return &Usercase{er}
}

func (u *Usercase) Equipment() equipment.Repository {
	return u.equipmentRepo
}
