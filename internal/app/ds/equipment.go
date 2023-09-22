package ds

type Equipment struct {
	ID           int `gorm:"primarykey"`
	Title        string
	Picture      string
	Description  string
	Status       string
	Count        int
	AvailableNow int
}
