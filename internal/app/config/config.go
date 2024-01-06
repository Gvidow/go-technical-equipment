package config

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/config"
)

type Config struct {
	ServiceHost string
	ServicePort string
	Mode        string

	JWT   JWTConfig
	Redis RedisConfig
}

func New(cfg *config.YAML) *Config {
	val := cfg.Get("app")
	c := &Config{
		Mode: val.Get("mode").String(),
	}
	val = val.Get("server")
	c.ServiceHost = val.Get("host").String()
	c.ServicePort = val.Get("port").String()

	c.JWT = JWTConfig{
		SecretToken: []byte("equipment"),
		Header:      "Authorization",
		TokenType:   "Bearer",
		SignMethod:  jwt.SigningMethodHS256,
		ExpiresIn:   3 * time.Hour,
	}

	c.Redis = RedisConfig{
		User:     os.Getenv(_envRedisUser),
		Password: os.Getenv(_envRedisPass),
		Host:     os.Getenv(_envRedisHost),
	}

	c.Redis.Port, _ = strconv.Atoi(os.Getenv(_envRedisPort))
	return c
}
