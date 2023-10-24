package request

import (
	"fmt"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &requestRepo{db}
}

type requestRepo struct {
	db *gorm.DB
}

func (r *requestRepo) GetLastEnteredRequestByUserID(userID int) (*ds.Request, error) {
	request := &ds.Request{}
	err := r.db.First(request, "creator = ?", userID).First(request, "status = ?", "entered").Error
	if err != nil {
		return nil, fmt.Errorf("get last entered request: %w", err)
	}
	return request, nil
}

func (r *requestRepo) SaveRequest(req *ds.Request) (*ds.Request, error) {
	err := r.db.Create(req).Error
	if err != nil {
		return nil, fmt.Errorf("save request: %w", err)
	}
	return req, nil
}
