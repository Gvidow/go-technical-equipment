package config

import "go.uber.org/config"

type Config struct {
	ServiceHost string
	ServicePort string
	Mode        string
}

func New(cfg *config.YAML) *Config {
	val := cfg.Get("app")
	c := &Config{
		Mode: val.Get("mode").String(),
	}
	val = val.Get("server")
	c.ServiceHost = val.Get("host").String()
	c.ServicePort = val.Get("port").String()
	return c
}
