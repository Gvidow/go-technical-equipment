package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretToken []byte
	Header      string
	TokenType   string
	SignMethod  jwt.SigningMethod
	ExpiresIn   time.Duration
}
