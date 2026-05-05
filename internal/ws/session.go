package ws

import (
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
)

type Session struct {
	ConnID string `json:"conn_id"`
	UserID int64  `json:"user_id"`
	SiteID int64  `json:"site_id"`
}

func GetSocketSession(connID string) *Session {
	r := redis.GetRedis()
	var sess Session
	ok, err := r.GetJSON(protocol.SessionKey(connID), &sess)
	if err != nil || !ok {
		return nil
	}
	return &sess
}
