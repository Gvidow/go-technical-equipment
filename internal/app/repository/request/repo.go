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

func (r *requestRepo) AddRequest(req *ds.Request) (*ds.Request, error) {
	err := r.db.Omit("moderator").Create(req).Error
	if err != nil {
		return nil, fmt.Errorf("save request: %w", err)
	}
	return req, nil
}

func (r *requestRepo) SaveRequest(req *ds.Request) error {
	req.ModeratorProfile = nil
	req.CreatorProfile = nil

	db := r.db
	if req.Moderator == 0 {
		db = db.Omit("moderator")
	}

	if err := db.Save(req).Error; err != nil {
		return fmt.Errorf("save request in storage: %w", err)
	}
	return nil
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
	return nil
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

func (r *requestRepo) UpdateRequest(requestID int, changes map[string]any) error {
	if err := r.db.Model(&ds.Request{ID: requestID}).Select("moderator", "creator").Updates(changes).Error; err != nil {
		return fmt.Errorf("request update %v error: %w", changes, err)
	}
	return nil
}

func (r *requestRepo) SaveUpdatedRequest(req *ds.Request) error {
	if err := r.db.Omit("formated_at", "completed_at").Save(req).Error; err != nil {
		return fmt.Errorf("save updated request: %w", err)
	}
	return nil
}

func (r *requestRepo) GetRequestWithFilter(cfg ds.FeedRequestConfig, userID int) ([]ds.Request, error) {
	feed := make([]ds.Request, 0)

	db := r.db.Where("status = ? OR status IN ? AND creator = ?", "operation", []string{"completed", "canceled"}, userID)

	if creator, ok := cfg.CreatorFilter(); ok {
		db = db.Where("creator = ?", creator)
	}
	if moderator, ok := cfg.ModeratorFilter(); ok {
		db = db.Where("moderator = ?", moderator)
	}

	// if creator, ok := cfg.CreatorProfileFilter(); ok {
	// 	db = db.InnerJoins("users")
	// }

	if status, ok := cfg.StatusFilter(); ok {
		db = db.Where("status = ?", status)
	}

	if createdAt, ok := cfg.CreatedAtFilter(); ok {
		db = db.Where("created_at = ?", createdAt)
	}
	if completedAt, ok := cfg.CreatedAtFilter(); ok {
		db = db.Where("completed_at = ?", completedAt)
	}
	if formatedAt, ok := cfg.CreatedAtFilter(); ok {
		db = db.Where("formated_at = ?", formatedAt)
	}
	if formatedAfter, ok := cfg.FormatedAfterFilter(); ok {
		db = db.Where("formated_at > ?", formatedAfter)
	}
	if formatedBefore, ok := cfg.FormatedBeforeFilter(); ok {
		db = db.Where("formated_at < ?", formatedBefore)
	}
	err := db.Preload("CreatorProfile").Preload("ModeratorProfile").Find(&feed).Error
	return feed, err
}

func (r *requestRepo) RevealEquipments(request *ds.Request) error {
	db := r.db.Raw("SELECT e.id, e.title, e.description, e.picture, e.status, o.count FROM request INNER JOIN orders o ON request.id = o.request_id INNER JOIN equipment e ON o.equipment_id = e.id WHERE request.id = ?;", request.ID)
	err := db.Scan(&request.Equipments).Error
	if err != nil {
		return fmt.Errorf("select equipments in request: %w", err)
	}
	return nil
}
