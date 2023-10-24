package order

import (
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
)

type Repository interface {
	OrderReplenishment(order ds.Order) error
}
