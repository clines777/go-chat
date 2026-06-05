package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"gochat/internal/infra"
)

var (
	redisConn *Client
	redisOnce sync.Once
)

type Client struct {
	Conn *goredis.Client
}

func GetRedis() *Client {
	redisOnce.Do(func() {
		cfg := infra.GetEnvConfig()
		redisConn = &Client{
			Conn: goredis.NewClient(&goredis.Options{
				Addr:         cfg.RedisHost,
				Password:     cfg.RedisPassword,
				DB:           cfg.RedisDB,
				DialTimeout:  cfg.RedisDialTimeout,
				ReadTimeout:  cfg.RedisDialTimeout,
				WriteTimeout: cfg.RedisDialTimeout,
				PoolSize:     20,
				MinIdleConns: 2,
			}),
		}
	})
	return redisConn
}

// Ping test Conn on startup
func (c *Client) Ping() error {
	if c == nil || c.Conn == nil {
		return fmt.Errorf("redis ws is nil")
	}
	return c.Conn.Ping(context.Background()).Err()
}

func (c *Client) GetString(key string) (string, bool, error) {
	val, err := c.Conn.Get(context.Background(), key).Result()
	if errors.Is(err, goredis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Client) SetString(key string, value string, exp time.Duration) error {
	return c.Conn.Set(context.Background(), key, value, exp).Err()
}

func (c *Client) Del(keys ...string) error {
	return c.Conn.Del(context.Background(), keys...).Err()
}

func (c *Client) Exists(key string) (bool, error) {
	n, err := c.Conn.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *Client) Expire(key string, exp time.Duration) error {
	return c.Conn.Expire(context.Background(), key, exp).Err()
}

func (c *Client) SetJSON(key string, value any, exp time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Conn.Set(context.Background(), key, b, exp).Err()
}

func (c *Client) GetJSON(key string, dest any) (bool, error) {
	val, err := c.Conn.Get(context.Background(), key).Bytes()
	if errors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(val, dest); err != nil {
		return false, err
	}
	return true, nil
}
