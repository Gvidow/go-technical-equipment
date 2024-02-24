package request

import "github.com/gvidow/go-technical-equipment/internal/app/ds"

type Repository interface {
	GetRequestByID(requestID int) (*ds.Request, error)
	GetEnteredRequest() (*ds.Request, error)
}
