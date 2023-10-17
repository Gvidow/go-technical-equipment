package equipment

import (
	"strings"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
)

type modelEquipment struct {
	data []ds.Equipment
}

func (m modelEquipment) GetByID(id int) (*ds.Equipment, error) {
	for _, equipment := range m.data {
		if equipment.ID == id {
			return &equipment, nil
		}
	}
	return nil, ErrEquipmentNotFound
}

func (m modelEquipment) SearchEquipmentsByTitle(title string) ([]ds.Equipment, error) {
	res := []ds.Equipment{}
	for _, equipment := range m.data {
		if strings.Contains(equipment.Title, title) {
			res = append(res, equipment)
		}
	}
	return res, nil
}

func (m modelEquipment) GetAllEquipments() ([]ds.Equipment, error) {
	return m.data, nil
}

func (m modelEquipment) DeleteEquipmentByID(id int) error {
	return nil
}
func (m modelEquipment) AddEquipment(*ds.Equipment) error {
	return nil
}

func NewStorageRepository() Repository {
	return modelEquipment{
		data: []ds.Equipment{
			{
				ID:          1,
				Title:       "Лазерная ручка",
				Picture:     "/upload/equipment/lizer.png",
				Description: "Помогает вести лекцию",
				Status:      "active",
				Count:       0,
				// AvailableNow: 0,
			},
			{
				ID:          2,
				Title:       "Проектор",
				Picture:     "/upload/equipment/projector.png",
				Description: "Очень полезное оборудование",
				Status:      "active",
				Count:       0,
				// AvailableNow: 0,
			},
			{
				ID:          3,
				Title:       "Экран",
				Picture:     "/upload/equipment/display.png",
				Description: "Уникальное оборудование",
				Status:      "active",
				Count:       0,
				// AvailableNow: 0,
			},
		},
	}
}
