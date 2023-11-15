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
	err := r.db.Where("status = 'entered'").First(request, "creator = ?", userID).Error
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

func (r *requestRepo) DeleteRequest(requestID int) error {
	db := r.db.Model(&ds.Request{}).Where("id = ?", requestID).Update("status", "deleted")
	if db.Error != nil {
		return fmt.Errorf("delete request from repository: %w", db.Error)
	}
	return nil
}

func (r *requestRepo) UpdateRequestStatus(requestID int, newStatus, oldStatusRequire string) error {
	db := r.db.Model(&ds.Request{}).Where("id = ?", requestID).Update("status", newStatus)
	if db.Error != nil {
		return fmt.Errorf("update status request on repository: %w", db.Error)
	}
	fmt.Println("ffdfdfDFd", db.Statement.SQL.String())
	return nil
}
