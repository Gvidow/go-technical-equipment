package equipment

import (
	"errors"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
)

var ErrEquipmentNotFound = errors.New("equipment not found")

type Repository interface {
	GetByID(id int) (*ds.Equipment, error)
	GetAllEquipments() ([]ds.Equipment, error)
	SearchEquipmentsByTitle(title string) ([]ds.Equipment, error)
	DeleteEquipmentByID(id int) error
	AddEquipment(eq *ds.Equipment) error
	ViewFeedEquipment(feedCfg ds.FeedEquipmentConfig) ([]ds.Equipment, error)
}
