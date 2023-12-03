package user

import (
	"fmt"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByID(userID int) (*ds.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db}
}

func (u *userRepo) GetUserByID(userID int) (*ds.User, error) {
	user := &ds.User{ID: userID}
	err := u.db.Find(user).Error
	if err != nil {
		return nil, fmt.Errorf("get user by id from storage: %w", err)
	}
	return user, nil
}
