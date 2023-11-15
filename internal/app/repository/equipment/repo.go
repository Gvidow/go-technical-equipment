package equipment

import (
	"fmt"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *equipmentRepo {
	repo := &equipmentRepo{db}
	return repo
}

type equipmentRepo struct {
	db *gorm.DB
}

func (er *equipmentRepo) GetAllEquipments() ([]ds.Equipment, error) {
	var equipments = make([]ds.Equipment, 0)
	db := er.db.Find(&equipments, "Status = ?", "active")
	return equipments, db.Error
}

func (er *equipmentRepo) GetByID(id int) (*ds.Equipment, error) {
	var equipment = &ds.Equipment{}
	db := er.db.First(equipment, id).First(equipment, "Status = ?", "active")
	return equipment, db.Error
}

func (er *equipmentRepo) SearchEquipmentsByTitle(title string) ([]ds.Equipment, error) {
	var equipments = make([]ds.Equipment, 0)
	db := er.db.Find(&equipments, "Status = ?", "active").Find(&equipments, "Title LIKE ?", "%"+title+"%")
	return equipments, db.Error
}

func (er *equipmentRepo) DeleteEquipmentByID(id int) error {
	db := er.db.Exec("UPDATE equipment SET status='delete' WHERE id = ?;", id)
	return db.Error
}

func (er *equipmentRepo) AddEquipment(eq *ds.Equipment) error {
	return er.db.Save(eq).Error
}

func (er *equipmentRepo) ViewFeedEquipment(cfg ds.FeedEquipmentConfig) ([]ds.Equipment, error) {
	feed := make([]ds.Equipment, 0)

	err := er.buildQueryFeedEquipment(cfg).Find(&feed).Error
	if err != nil {
		return nil, fmt.Errorf("get feed equipment from storage: %w", err)
	}
	return feed, nil
}

func (er *equipmentRepo) buildQueryFeedEquipment(cfg ds.FeedEquipmentConfig) *gorm.DB {
	db := er.db

	if cfg.InStock {
		db = db.Where("count > 0")
	}

	if title, ok := cfg.TitleFilter(); ok {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}

	switch cfg.Status {
	case ds.Active:
		db = db.Where("status = 'active'")
	case ds.Delete:
		db = db.Where("status = 'delete'")
	}

	if date, ok := cfg.DateCreateFilter(); ok {
		db = db.Where("created_at > ?", date.Format("01-02-2006"))
	}
	return db
}
