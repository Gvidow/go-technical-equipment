package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gvidow/go-technical-equipment/internal/app/config"
	redis "github.com/redis/go-redis/v9"
)

const _servicePrefix = "equipment_service."

type Client struct {
	cfg    config.RedisConfig
	client *redis.Client
}

func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	client := &Client{
		cfg: cfg,
		client: redis.NewClient(&redis.Options{
			Password:    cfg.Password,
			Username:    cfg.User,
			Addr:        cfg.Host + ":" + strconv.Itoa(cfg.Port),
			DB:          0,
			DialTimeout: cfg.DialTimeout,
			ReadTimeout: cfg.ReadTimeout,
		}),
	}
	if _, err := client.client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
