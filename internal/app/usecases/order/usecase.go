package order

import (
	"errors"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/equipment"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/order"
)

var ErrEquipmentNotFound = errors.New("not found equipment")

type Usecase struct {
	repo   order.Repository
	repoEq equipment.Repository
}

func NewUsecase(repo order.Repository, eq equipment.Repository) *Usecase {
	return &Usecase{repo, eq}
}

func (u *Usecase) AddEquipmentInRequest(equipmentID, requestID int) error {
	eq, err := u.repoEq.GetByID(equipmentID)
	if err != nil || eq.Status == "deleted" {
		return ErrEquipmentNotFound
	}

	ord := ds.Order{
		EquipmentID: equipmentID,
		RequestID:   requestID,
	}

	err = u.repo.OrderReplenishment(ord)

	if err == order.ErrEquipmentAlreadyAdd {
		return u.repo.OrderInc(ord)
	}
	return err
}

func (u *Usecase) DeleteEquipmentFromRequest(equipmentID, requestID int) error {
	return u.repo.DropOrder(ds.Order{EquipmentID: equipmentID, RequestID: requestID})
}

func (u *Usecase) EditCountEquipmentsInRequest(equipmentID, requestID, newCount int) error {
	if newCount == 0 {
		return u.repo.DropOrder(ds.Order{EquipmentID: equipmentID, RequestID: requestID})
	}

	return u.repo.EditCountEquipmentInOrder(ds.Order{EquipmentID: equipmentID, RequestID: requestID}, newCount)
}
