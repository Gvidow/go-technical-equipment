package request

import "github.com/gvidow/go-technical-equipment/internal/app/ds"

type Repository interface {
	GetRequestWithFilter(cfg ds.FeedRequestConfig, userID int) ([]ds.Request, error)
	GetRequestByID(requestID int) (*ds.Request, error)
	GetLastEnteredRequestByUserID(userID int) (*ds.Request, error)
	AddRequest(req *ds.Request) (*ds.Request, error)
	SaveRequest(req *ds.Request) error
	DeleteRequest(requestID int) error
	UpdateRequestStatus(requestID int, newStatus, oldStatusRequire string) error
	SaveUpdatedRequest(req *ds.Request) error
	UpdateRequest(requestID int, changes map[string]any) error
	RevealEquipments(request *ds.Request) error
	UpdateReverted(reqID int, reverted bool) error
}
