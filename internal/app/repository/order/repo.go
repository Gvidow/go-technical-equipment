package order

import (
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &orderRepo{db}
}

type orderRepo struct {
	db *gorm.DB
}

func (o *orderRepo) OrderReplenishment(order ds.Order) error {
	return o.db.Create(order).Error
}
