package order

import (
	"errors"
	"fmt"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrEquipmentAlreadyAdd = errors.New("equipment already add")

func NewRepository(db *gorm.DB) Repository {
	return &orderRepo{db}
}

type orderRepo struct {
	db *gorm.DB
}

func (o *orderRepo) OrderReplenishment(order ds.Order) error {
	// equipment := &ds.Equipment{}
	// err := o.db.First(equipment, "id = ?", order.EquipmentID).Error
	// if err != nil {
	// 	return fmt.Errorf("get equipment for adding for replenishment order: %w", err)
	// }

	// err = o.db.Model(&ds.Request{ID: order.RequestID}).Association("Equipments").Append(equipment)
	// if err != nil {
	// 	return fmt.Errorf("order replenishment: %w", err)
	// }

	err := o.db.Create(order).Error
	t := &pgconn.PgError{}
	ok := errors.As(err, &t)
	if ok && t.Code == "23505" {
		return ErrEquipmentAlreadyAdd
	}
	return err
}

func (o *orderRepo) DropOrder(order ds.Order) error {
	err := o.db.Model(&ds.Request{ID: order.RequestID}).Association("Equipments").Delete(&ds.Equipment{ID: order.EquipmentID})
	if err != nil {
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
