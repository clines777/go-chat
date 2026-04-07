package redis

import (
	"context"
	"encoding/json"
	"fmt"
	goredis "github.com/redis/go-redis/v9"
	"gochat/internal/infra"
	"time"
)

var redisConn *Client

type Client struct {
	Conn *goredis.Client
	ctx  context.Context
}

func GetRedis() *Client {
	if redisConn == nil {
		cfg := infra.GetEnvConfig()
		redisConn = &Client{}
		redisConn.Conn = goredis.NewClient(&goredis.Options{
			Addr:         cfg.RedisHost,
			Password:     cfg.RedisPassword,
			DB:           cfg.RedisDB,
			DialTimeout:  cfg.RedisDialTimeout,
			ReadTimeout:  cfg.RedisDialTimeout,
			WriteTimeout: cfg.RedisDialTimeout,
			PoolSize:     20,
			MinIdleConns: 2,
		})
	}

	redisConn.ctx = context.Background()

	return redisConn
}

func GetRedisWithContext(ctx context.Context) *Client {
	if redisConn == nil {
		cfg := infra.GetEnvConfig()
		redisConn = &Client{}
		redisConn.Conn = goredis.NewClient(&goredis.Options{
			Addr:         cfg.RedisHost,
			Password:     cfg.RedisPassword,
			DB:           cfg.RedisDB,
			DialTimeout:  cfg.RedisDialTimeout,
			ReadTimeout:  cfg.RedisDialTimeout,
			WriteTimeout: cfg.RedisDialTimeout,
			PoolSize:     20,
			MinIdleConns: 2,
		})
	}
	redisConn.ctx = ctx

	return redisConn
}

// Ping test Conn on startup
func (c *Client) Ping() error {
	if c == nil || c.Conn == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.Conn.Ping(c.ctx).Err()
}

func (c *Client) GetString(key string) (string, bool, error) {
	val, err := c.Conn.Get(c.ctx, key).Result()
	if err == goredis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Client) SetString(key string, value string, exp time.Duration) error {
	return c.Conn.Set(c.ctx, key, value, exp).Err()
}

func (c *Client) Del(keys ...string) error {
	return c.Conn.Del(c.ctx, keys...).Err()
}

func (c *Client) Exists(key string) (bool, error) {
	n, err := c.Conn.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *Client) Expire(key string, exp time.Duration) error {
	return c.Conn.Expire(c.ctx, key, exp).Err()
}

func (c *Client) SetJSON(key string, value any, exp time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Conn.Set(c.ctx, key, b, exp).Err()
}

func (c *Client) GetJSON(key string, dest any) (bool, error) {
	val, err := c.Conn.Get(c.ctx, key).Bytes()
	if err == goredis.Nil {
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
