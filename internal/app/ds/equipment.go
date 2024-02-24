package ds

type Equipment struct {
	ID          int       `json:"id" gorm:"primarykey"`
	Title       string    `json:"title"`
	Picture     string    `json:"picture"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Count       int       `json:"count,omitempty"`
	Requests    []Request `json:"-" gorm:"many2many:orders;"`
}
