package ds

import (
	"time"
)

type Request struct {
	ID               int         `json:"id" gorm:"primary_key"`
	Status           string      `json:"status"`
	Moderator        int         `json:"-"`
	Creator          int         `json:"-"`
	CreatedAt        *time.Time  `gorm:"created_at;" json:"created_at"`
	FormatedAt       *time.Time  `gorm:"formated_at;null" json:"formated_at"`
	CompletedAt      *time.Time  `gorm:"completed_at;null" json:"completed_at"`
	CreatorProfile   *User       `gorm:"foreignKey:creator;references:id" json:"creator_profile"`
	ModeratorProfile *User       `gorm:"foreignKey:moderator;references:id" json:"moderator_profile"`
	Equipments       []Equipment `gorm:"many2many:orders;" json:"equipments,omitempty"`
	Reverted         *bool       `json:"reverted,omitempty"`
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
