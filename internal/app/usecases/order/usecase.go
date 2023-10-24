package order

import (
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/order"
)

type Usecase struct {
	repo order.Repository
}

func NewUsecase(repo order.Repository) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) AddEquipmentInRequest(equipmentID, requestID int) error {
	return u.repo.OrderReplenishment(ds.Order{
		EquipmentID: equipmentID,
		RequestID:   requestID,
	})
}
