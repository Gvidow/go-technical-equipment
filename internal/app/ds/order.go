package ds

type Order struct {
	EquipmentID int `gorm:"equipment_id;primary_key"`
	RequestID   int `gorm:"request_id;primary_key"`
}

func (o Order) TableName() string {
	return "orders"
}
