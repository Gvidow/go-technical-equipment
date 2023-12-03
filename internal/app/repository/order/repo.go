package order

import (
	"fmt"

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

func (o *orderRepo) DropOrder(order ds.Order) error {
	if err := o.db.Delete(order).Error; err != nil {
		return fmt.Errorf("drop order from storage: %w", err)
	}
	return nil
}

func (o *orderRepo) OrderInc(order ds.Order) error {
	if err := o.editCount(order, gorm.Expr("count + ?", 1)); err != nil {
		return fmt.Errorf("order increment error: %w", err)
	}
	return nil
}

func (o *orderRepo) OrderDec(order ds.Order) error {
	if err := o.editCount(order, gorm.Expr("count - ?", 1)); err != nil {
		return fmt.Errorf("order decrement error: %w", err)
	}
	return nil
}

func (o *orderRepo) EditCountEquipmentInOrder(order ds.Order, newCount int) error {
	if err := o.editCount(order, newCount); err != nil {
		return fmt.Errorf("order update count on %d: %w", newCount, err)
	}
	return nil
}

func (o *orderRepo) editCount(order ds.Order, value any) error {
	return o.db.Model(order).Where("equipment_id = ?", order.EquipmentID).Where("request_id = ?", order.RequestID).
		UpdateColumn("count", value).Error
}
