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

func (r *requestRepo) GetEnteredRequest() (*ds.Request, error) {
	request := &ds.Request{}
	err := r.db.Where("status = 'entered'").First(request).Error
	if err != nil {
		return nil, fmt.Errorf("get entered request: %w", err)
	}
	return request, nil
}

func (r *requestRepo) GetRequestByID(requestID int) (*ds.Request, error) {
	req := &ds.Request{ID: requestID}
	if err := r.db.
		Preload("CreatorProfile").
		Preload("ModeratorProfile").
		Preload("Equipments").
		First(req).Error; err != nil {
		return nil, fmt.Errorf("get request by id from storage: %w", err)
	}
	return req, nil
}
