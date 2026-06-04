package session

import (
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
	"time"
)

type Session struct {
	Scene     protocol.Scene `json:"scene"`
	InGroupID int            `json:"in_group_id,omitempty"`
	UserID    int            `json:"user_id"`
	ConnID    string         `json:"conn_id"`
	Username  string         `json:"username"`
}

func Get(userID int, connID string) *Session {
	r := redis.GetRedis()
	var sess Session
	ok, err := r.GetJSON(protocol.SessionKey(userID, connID), &sess)
	if err != nil || !ok {
		return nil
	}
	return &sess
}

func Set(userID int, connID string, sess *Session) error {
	r := redis.GetRedis()
	if err := r.SetJSON(protocol.SessionKey(userID, connID), sess, 24*time.Hour); err != nil {
		return err
	}
	return nil
}
