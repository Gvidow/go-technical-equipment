package ds

import "time"

type Equipment struct {
	ID          int `gorm:"primarykey"`
	Title       string
	Picture     string
	Description string
	Status      string
	Count       int
}

type equipmentStatus uint

const (
	_ equipmentStatus = iota
	Active
	Delete
	All
)

type FeedEquipmentConfig struct {
	createdAfter time.Time       `schema:"createdAfter"`
	titleFilter  string          `schema:"title"`
	Status       equipmentStatus `schema:"status"`
	InStock      bool            `schema:"inStock"`

	hasFilterDateCreate bool
	hasFilterTitle      bool
}

func (f *FeedEquipmentConfig) TitleFilter() (string, bool) {
	return f.titleFilter, f.hasFilterTitle
}

func (f *FeedEquipmentConfig) DateCreateFilter() (time.Time, bool) {
	return f.createdAfter, f.hasFilterDateCreate
}

func (f *FeedEquipmentConfig) SetTitleFilter(title string) {
	f.titleFilter, f.hasFilterTitle = title, true
}

func (f *FeedEquipmentConfig) SetDateCreateFilter(date time.Time) {
	f.createdAfter, f.hasFilterDateCreate = date, true
}
