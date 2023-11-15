package request

import (
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/request"
)

type Usecase struct {
	repo request.Repository
}

func NewUsecase(repo request.Repository) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) GettingUserLastRequest(userID int) (*ds.Request, error) {
	return u.repo.GetLastEnteredRequestByUserID(userID)
}

func (u *Usecase) CreateDraftRequest(userID int) (*ds.Request, error) {
	request := &ds.Request{
		Creator:   userID,
		Moderator: userID,
		Status:    "entered",
	}
	return u.repo.SaveRequest(request)
}

func (u *Usecase) DropRequest(requestID int) error {
	return u.repo.DeleteRequest(requestID)
}

func (u *Usecase) ChangeStatusRequest(requestID int, newStatus, oldStatusRequire string) error {
	return u.repo.UpdateRequestStatus(requestID, newStatus, oldStatusRequire)
}
