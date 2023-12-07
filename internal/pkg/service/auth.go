package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/usecases/auth"
)

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (s *Service) Login(c *gin.Context) {
	req := &loginReq{}
	err := json.NewDecoder(c.Request.Body).Decode(req)
	defer c.Request.Body.Close()
	if err != nil {
		s.log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось распарсить тело запроса"})
		return
	}

	token, err := s.authCase.Login(req.Login, req.Password, s.cfg.JWT)
	if err == auth.ErrIncorrectCredentials {
		s.log.Info(err.Error())
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "неправильные логин или пароль"})
		return
	}

	if err != nil {
		s.log.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ошибка при входе в систему"})
		return
	}

	c.JSON(http.StatusCreated, loginResp{
		ExpiresIn:   int(s.cfg.JWT.ExpiresIn),
		AccessToken: token,
		TokenType:   s.cfg.JWT.TokenType,
	})
}

func (s *Service) Signup(c *gin.Context) {
	cred := &ds.Credentials{}
	err := json.NewDecoder(c.Request.Body).Decode(cred)
	defer c.Request.Body.Close()
	if err != nil {
		s.log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "не удалось распарсить тело запроса"})
		return
	}

	if err = s.authCase.Signup(cred); err != nil {
		s.log.Warn(err.Error())
		switch err {
		case auth.ErrBadCredentials:
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "неверные регистрационные данные"})
		default:
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "регистрация не прошла"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "регистрация прошла успешно"})
}

func (s *Service) Logout(c *gin.Context) {
	err := s.authCase.Logout(c, c.GetHeader(s.cfg.JWT.Header), s.cfg.JWT)
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "вы успешно вышли из аккаунта"})
	case auth.ErrInvalidToken:
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "недействительный токен"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ну удалось вас разлогинить"})
	}
}
