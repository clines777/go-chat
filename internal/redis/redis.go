package redis

import (
	"context"
	"fmt"
	goredis "github.com/redis/go-redis/v9"
	"gochat/internal/config"
)

type Client struct {
	Rdb *goredis.Client
}

func Init(cfg *config.EnvConfig) (*Client, error) {

	rdb := goredis.NewClient(&goredis.Options{
		Addr:         cfg.RedisHost,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  cfg.RedisDialTimeout,
		ReadTimeout:  cfg.RedisDialTimeout,
		WriteTimeout: cfg.RedisDialTimeout,
		PoolSize:     20,
		MinIdleConns: 2,
	})

	return &Client{Rdb: rdb}, nil
}

func (c *Client) Close() error {
	if c == nil || c.Rdb == nil {
		return nil
	}
	return c.Rdb.Close()
}

// Ping test conn on startup
func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.Rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.Rdb.Ping(ctx).Err()
}
