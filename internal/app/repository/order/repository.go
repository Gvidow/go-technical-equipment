package order

import (
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
)

type Repository interface {
	OrderReplenishment(order ds.Order) error
	DropOrder(order ds.Order) error
	OrderInc(order ds.Order) error
	OrderDec(order ds.Order) error
	EditCountEquipmentInOrder(order ds.Order, newCount int) error
}
