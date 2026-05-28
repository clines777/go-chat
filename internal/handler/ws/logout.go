package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

type logoutReq struct {
	ResumeToken string `json:"resume_token"`
}

func Logout(ctx *ws.Ctx) *protocol.Payload {
	var req logoutReq
	_ = json.Unmarshal(ctx.Payload.Data, &req)

	sess := session.Get(ctx.Client.ConnID)

	if req.ResumeToken != "" {
		_ = redis.GetRedis().Del(protocol.ResumeTokenKey(req.ResumeToken))
	}

	log.Printf("[Logout] OK user=%d", sess.UserID)
	return &protocol.Payload{MsgType: protocol.LogoutOK, Meta: ctx.Payload.Meta}
}
