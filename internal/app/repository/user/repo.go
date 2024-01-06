package user

import (
	"errors"
	"fmt"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")

type Repository interface {
	GetUserByID(userID int) (*ds.User, error)
	AddUser(user *ds.User) (*ds.User, error)
	GetUserByUsernameOrEmail(login string) (*ds.User, error)
	GetUserByUsername(username string) (*ds.User, error)
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

func (u *userRepo) AddUser(user *ds.User) (*ds.User, error) {
	err := u.db.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("add new user %w", err)
	}
	return user, nil
}

func (u *userRepo) GetUserByUsernameOrEmail(login string) (*ds.User, error) {
	user := &ds.User{}
	if err := u.db.Where("email = ?", login).Or("username = ?", login).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("get user by username or email from storage: %w", err)
	}
	return user, nil
}

func (u *userRepo) GetUserByUsername(username string) (*ds.User, error) {
	user := &ds.User{
		Username: username,
	}
	if err := u.db.Where(user).First(user).Error; err != nil {
		return nil, fmt.Errorf("get user by username from storage: %w", err)
	}
	return user, nil
}
