package ds

import (
	"strconv"
	"time"
)

var _emptyTime = time.Time{}

type Request struct {
	ID               int `gorm:"primary_key"`
	Status           string
	Moderator        int         `json:"moderator,omitempty"`
	Creator          int         `json:"creator,omitempty"`
	CreatedAt        *time.Time  `gorm:"created_at;" json:",omitempty"`
	FormatedAt       *time.Time  `gorm:"formated_at;null" json:",omitempty"`
	CompletedAt      *time.Time  `gorm:"completed_at;null" json:",omitempty"`
	CreatorProfile   *User       `gorm:"-" json:",omitempty"`
	ModeratorProfile *User       `gorm:"-" json:",omitempty"`
	Equipments       []Equipment `gorm:"-" json:",omitempty"`
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

type FeedRequestConfig struct {
	creator      int
	moderator    int
	status       string
	created_at   time.Time
	formated_at  time.Time
	completed_at time.Time
}

func (f *FeedRequestConfig) SetCreatorFilter(creator string) error {
	id, err := strconv.ParseInt(creator, 10, 64)
	if err != nil {
		return err
	}
	f.creator = int(id)
	return nil
}

func (f *FeedRequestConfig) SetModeratorFilter(moderator string) error {
	id, err := strconv.ParseInt(moderator, 10, 64)
	if err != nil {
		return err
	}
	f.moderator = int(id)
	return nil
}

func (f *FeedRequestConfig) SetStatusFilter(status string) {
	f.status = status
}

func (f *FeedRequestConfig) SetCreatedFilter(date string) error {
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	f.created_at = t
	return nil
}

func (f *FeedRequestConfig) SetCompletedFilter(date string) error {
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	f.completed_at = t
	return nil
}

func (f *FeedRequestConfig) SetFormatedFilter(date string) error {
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	f.formated_at = t
	return nil
}

func (f *FeedRequestConfig) CreatorFilter() (int, bool) {
	if f.creator == 0 {
		return 0, false
	}
	return f.creator, true
}

func (f *FeedRequestConfig) ModeratorFilter() (int, bool) {
	if f.moderator == 0 {
		return 0, false
	}
	return f.moderator, true
}

func (f *FeedRequestConfig) StatusFilter() (string, bool) {
	if f.status == "" {
		return "", false
	}
	return f.status, true
}

func (f *FeedRequestConfig) CreatedAtFilter() (time.Time, bool) {
	if f.created_at == _emptyTime {
		return _emptyTime, false
	}
	return f.created_at, true
}

func (f *FeedRequestConfig) FormatedAtFilter() (time.Time, bool) {
	if f.formated_at == _emptyTime {
		return _emptyTime, false
	}
	return f.formated_at, true
}
func (f *FeedRequestConfig) CompletedAtFilter() (time.Time, bool) {
	if f.completed_at == _emptyTime {
		return _emptyTime, false
	}
	return f.completed_at, true
}
