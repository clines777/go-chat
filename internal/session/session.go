package session

import (
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
	"time"
)

type Session struct {
	Scene     protocol.Scene `json:"scene"`
	InGroupId int32          `json:"in_group_id,omitempty"`
	UserID    int64          `json:"user_id"`
	ConnID    string         `json:"conn_id"`
	SiteBid   string         `json:"site_bid"`
}

func Get(connID string) *Session {
	r := redis.GetRedis()
	var sess Session
	ok, err := r.GetJSON(protocol.SessionKey(connID), &sess)
	if err != nil || !ok {
		return nil
	}
	return &sess
}

func Set(connID string, sess *Session) error {
	r := redis.GetRedis()
	if err := r.SetJSON(protocol.SessionKey(connID), sess, 24*time.Hour); err != nil {
		return err
	}

	return nil
}
