package equipment

import (
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
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
