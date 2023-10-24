package ds

import "time"

type Request struct {
	ID          int `gorm:"primary_key"`
	Status      string
	Moderator   int
	Creator     int
	FormatedAt  time.Time `gorm:"formated_at,ommitempty"`
	CompletedAt time.Time `gorm:"completed_at,ommitempty"`
}

func (r *Request) Id() int {
	if r == nil {
		return 0
	}
	return r.ID
}

func (r *Request) TableName() string {
	return "request"
}
