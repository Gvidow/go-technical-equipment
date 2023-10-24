package request

import "github.com/gvidow/go-technical-equipment/internal/app/ds"

type Repository interface {
	GetLastEnteredRequestByUserID(userID int) (*ds.Request, error)
	SaveRequest(req *ds.Request) (*ds.Request, error)
}
