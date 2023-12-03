package request

import (
	"fmt"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/request"
	"github.com/gvidow/go-technical-equipment/internal/app/repository/user"
)

var _ok = struct{}{}

var _permissionsChangeStatuses = map[ds.Role]map[string]struct{}{
	ds.RegularUser: {
		"operation": _ok,
		"deleted":   _ok,
	},
	ds.Moderator: {
		"completed": _ok,
		"canceled":  _ok,
	},
}

var (
	ErrNotAccess         = fmt.Errorf("not access")
	ErrRoleHaveNotAccess = fmt.Errorf("the role does not have access")
)

type Usecase struct {
	repo     request.Repository
	userRepo user.Repository
}

func NewUsecase(repo request.Repository, userRepo user.Repository) *Usecase {
	return &Usecase{repo, userRepo}
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

func (u *Usecase) ChangeStatusRequest(userID, requestID int, newStatus string, requestedRole ds.Role) error {
	req, err := u.repo.GetRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("change status request from %s: %w", requestedRole, err)
	}

	if _, ok := _permissionsChangeStatuses[requestedRole][newStatus]; !ok {
		return ErrRoleHaveNotAccess
	}

	if requestedRole == ds.RegularUser && req.Creator != userID {
		return ErrNotAccess
	}

	req.Status = newStatus

	err = u.repo.SaveUpdatedRequest(req)
	if err != nil {
		return fmt.Errorf("change status request from role %s: %w", requestedRole, err)
	}
	return nil
}

func (u *Usecase) GetRequestByID(requestID int) (*ds.Request, error) {
	request, err := u.repo.GetRequestByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("get request by id: %w", err)
	}

	err = u.revealCreator(request)
	if err != nil {
		return nil, err
	}

	err = u.revealModerator(request)
	if err != nil {
		return nil, err
	}

	err = u.repo.RevealEquipments(request)
	if err != nil {
		return nil, fmt.Errorf("reveal equipments: %w", err)
	}

	return request, nil
}

func (u *Usecase) EditRequest(requestID int, changes map[string]any) error {
	return u.repo.UpdateRequest(requestID, changes)
}

func (u *Usecase) GetFeedRequests(cfg ds.FeedRequestConfig) ([]ds.Request, error) {
	return u.repo.GetRequestWithFilter(cfg)
}

func (u *Usecase) revealCreator(request *ds.Request) error {
	user, err := u.userRepo.GetUserByID(request.Creator)
	if err != nil {
		return fmt.Errorf("reveal creator: %w", err)
	}
	request.Creator = 0
	request.CreatorProfile = user
	return nil
}

func (u *Usecase) revealModerator(request *ds.Request) error {
	user, err := u.userRepo.GetUserByID(request.Moderator)
	if err != nil {
		return fmt.Errorf("reveal moderator: %w", err)
	}
	request.Moderator = 0
	request.ModeratorProfile = user
	return nil
}
