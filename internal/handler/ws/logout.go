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

// Logout - 前端頁面加了一個登出鈕, 方便測試時換帳號用.
func Logout(ctx *ws.Ctx) *protocol.Payload {
	var req logoutReq
	_ = json.Unmarshal(ctx.Payload.Data, &req)

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)

	if req.ResumeToken != "" {
		r := redis.GetRedis()
		_ = r.Del(protocol.SessionKey(ctx.Client.UserId, ctx.Client.ConnID))
		_ = r.Del(protocol.ResumeTokenKey(req.ResumeToken))
		if ctx.Client.ApiToken != "" {
			_ = r.Del(protocol.ApiTokenKey(ctx.Client.ApiToken))
		}
	}

	log.Printf("[Logout] OK user=%d", sess.UserID)
	return &protocol.Payload{MsgType: protocol.LogoutOK, Meta: ctx.Payload.Meta}
}
