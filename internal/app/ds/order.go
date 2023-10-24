package ds

type Order struct {
	EquipmentID int `gorm:"equipment_id"`
	RequestID   int `gorm:"request_id`
}

func (o Order) TableName() string {
	return "orders"
}
