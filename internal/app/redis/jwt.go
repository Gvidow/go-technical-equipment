package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const _jwtPrefix = "jwt."

func getJWTKey(token string) string {
	return _servicePrefix + _jwtPrefix + token
}

func (c *Client) WriteJWTToBlackList(ctx context.Context, jwtStr string, jwtTTL time.Duration) error {
	return c.client.Set(ctx, getJWTKey(jwtStr), true, jwtTTL).Err()
}

func (c *Client) CheckJWTInBlackList(ctx context.Context, jwtStr string) (bool, error) {
	err := c.client.Get(ctx, getJWTKey(jwtStr)).Err()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("check jwt in black list: %w", err)
	}
	return err == redis.Nil, nil
}
