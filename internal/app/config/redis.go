package config

import "time"

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
	_envRedisHost = "REDIS_HOST"
	_envRedisPort = "REDIS_PORT"
	_envRedisUser = "REDIS_USER"
	_envRedisPass = "REDIS_PASSWORD"
)
