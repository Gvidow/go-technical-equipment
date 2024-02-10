package ds

import (
	"strconv"
	"time"
)

var _emptyTime = time.Time{}

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

type filter[T any] struct {
	val T
	has bool
}

func (f *filter[T]) SetFilter(val T) { f.val, f.has = val, true }

func (f *filter[T]) Filter() (T, bool) { return f.val, f.has }

func (f *filter[T]) Has() bool { return f.has }

func (f *filter[T]) Clean() { f.val, f.has = *new(T), false }

type FeedRequestConfig struct {
	creator          int
	moderator        int
	creatorProfile   filter[string]
	moderatorProfile filter[string]
	status           string
	createdAt        time.Time
	formatedAt       time.Time
	formatedAfter    time.Time
	formatedBefore   time.Time
	completedAt      time.Time
}

func (f *FeedRequestConfig) SetCreatorFilterInt(creator int) {
	f.creator = creator
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
	f.createdAt = t
	return nil
}

func (f *FeedRequestConfig) SetCompletedFilter(date string) error {
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	f.completedAt = t
	return nil
}

func (f *FeedRequestConfig) SetFormatedFilter(date string) error {
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	f.formatedAt = t
	return nil
}

func (f *FeedRequestConfig) SetCreatorProfileFilter(username string) {
	f.creatorProfile.SetFilter(username)
}

func (f *FeedRequestConfig) SetModeratorProfileFilter(username string) {
	f.moderatorProfile.SetFilter(username)
}

func (f *FeedRequestConfig) CleanCreatorProfileFilter() {
	f.creatorProfile.Clean()
}

func (f *FeedRequestConfig) CleanModeratorProfileFilter() {
	f.moderatorProfile.Clean()
}

func (f *FeedRequestConfig) CreatorProfileFilter() (string, bool) {
	return f.creatorProfile.Filter()
}

func (f *FeedRequestConfig) ModeratorProfileFilter() (string, bool) {
	return f.moderatorProfile.Filter()
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
	if f.createdAt == _emptyTime {
		return _emptyTime, false
	}
	return f.createdAt, true
}

func (f *FeedRequestConfig) FormatedAtFilter() (time.Time, bool) {
	if f.formatedAt == _emptyTime {
		return _emptyTime, false
	}
	return f.formatedAt, true
}
func (f *FeedRequestConfig) CompletedAtFilter() (time.Time, bool) {
	if f.completedAt == _emptyTime {
		return _emptyTime, false
	}
	return f.completedAt, true
}

func (f *FeedRequestConfig) SetFormatedAfter(date string) error {
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	f.formatedAfter = t
	return nil
}

func (f *FeedRequestConfig) SetFormatedBefore(date string) error {
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return err
	}
	f.formatedBefore = t
	return nil
}

func (f *FeedRequestConfig) FormatedAfterFilter() (time.Time, bool) {
	if f.formatedAfter == _emptyTime {
		return _emptyTime, false
	}
	return f.formatedAfter, true
}

func (f *FeedRequestConfig) FormatedBeforeFilter() (time.Time, bool) {
	if f.formatedBefore == _emptyTime {
		return _emptyTime, false
	}
	return f.formatedBefore, true
}
