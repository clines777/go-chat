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

// Logout - 方便測試時換帳號用.
func Logout(ctx *ws.Ctx) *protocol.Payload {
	var req logoutReq
	_ = json.Unmarshal(ctx.Payload.Data, &req)

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)

	if req.ResumeToken != "" {
		r := redis.GetRedis()
		_ = r.Del(protocol.SessionKey(ctx.Client.UserId, ctx.Client.ConnID))
		_ = r.Del(protocol.ResumeTokenKey(req.ResumeToken))
		if sess != nil && sess.ApiToken != "" {
			_ = r.Del(protocol.ApiTokenKey(sess.ApiToken))
		}
	}

	log.Printf("[Logout] OK user=%d", sess.UserID)
	return &protocol.Payload{MsgType: protocol.LogoutOK, Meta: ctx.Payload.Meta}
}
