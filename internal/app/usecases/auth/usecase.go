package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/redis"
	userRepo "github.com/gvidow/go-technical-equipment/internal/app/repository/user"
	"github.com/gvidow/go-technical-equipment/pkg/crypto"
)

const _maxLenCred = 50

var (
	ErrBadCredentials       = errors.New("bad credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrIncorrectCredentials = errors.New("incorrect credentials")
)

type Usecase struct {
	repo      userRepo.Repository
	blackList *redis.Client
}

func NewUsecase(repo userRepo.Repository, client *redis.Client) *Usecase {
	return &Usecase{repo, client}
}

func (u *Usecase) Login(login, password string, cfg config.JWTConfig) (string, error) {
	user, err := u.repo.GetUserByUsernameOrEmail(login)
	if err != nil && err != userRepo.ErrRecordNotFound {
		return "", fmt.Errorf("login: %w", err)
	}

	if err == userRepo.ErrRecordNotFound || !isCorrectPassword(user.Password, password) {
		return "", ErrIncorrectCredentials
	}

	token := jwt.NewWithClaims(cfg.SignMethod, &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(cfg.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		UserID: user.ID,
		Role:   user.GetRole(),
	})

	strToken, err := token.SignedString(cfg.SecretToken)
	if err != nil {
		return "", fmt.Errorf("signed string token for login: %w", err)
	}
	return strToken, nil
}

func (u *Usecase) Signup(cred *ds.Credentials) error {
	if !checkCred(cred.Username) || !checkCred(cred.Email) || !checkCred(cred.Password) {
		return ErrBadCredentials
	}

	user := &ds.User{
		Username: cred.Username,
		Email:    cred.Email,
	}

	salt := crypto.NewSalt()
	user.Password = salt + crypto.Hash(cred.Password, salt)

	user.SetRole(ds.RegularUser)

	if _, err := u.repo.AddUser(user); err != nil {
		return fmt.Errorf("signup user: %w", err)
	}
	return nil
}

func (u *Usecase) Logout(ctx context.Context, token string, cfg config.JWTConfig) error {
	if !strings.HasPrefix(token, cfg.TokenType+" ") {
		return ErrInvalidToken
	}
	token = token[len(cfg.TokenType)+1:]

	if ok, err := u.blackList.CheckJWTInBlackList(ctx, token); err != nil && !ok {
		return fmt.Errorf("check jwt in black list for logout: %w", err)
	} else if !ok {
		return ErrInvalidToken
	}

	if err := u.blackList.WriteJWTToBlackList(ctx, token, cfg.ExpiresIn); err != nil {
		return fmt.Errorf("write to blacklist for logout: %w", err)
	}
	return nil
}

func checkCred(row string) bool {
	if len(row) == 0 || len(row) > _maxLenCred {
		return false
	}
	return true
}

func isCorrectPassword(hash, pass string) bool {
	if len(hash) < crypto.LenSalt {
		return false
	}

	salt := hash[:crypto.LenSalt]
	return hash[crypto.LenSalt:] == crypto.Hash(pass, salt)
}
