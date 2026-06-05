package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

type logoutReq struct {
	ResumeToken string `json:"resume_token"`
}

// Logout - 前端頁面加了一個登出鈕, 方便測試時換帳號用.
func Logout(ctx *ws.Ctx) *protocol.Payload {
	var req logoutReq
	_ = json.Unmarshal(ctx.Payload.Data, &req)

	sess := ctx.Session

	r := redis.GetRedis()
	_ = r.Del(protocol.SessionKey(ctx.Client.UserId, ctx.Client.ConnID))
	if ctx.Client.ApiToken != "" {
		_ = r.Del(protocol.ApiTokenKey(ctx.Client.ApiToken))
	}
	if req.ResumeToken != "" {
		_ = r.Del(protocol.ResumeTokenKey(req.ResumeToken))
	}

	log.Printf("[Logout] OK user=%d", sess.UserID)
	return &protocol.Payload{MsgType: protocol.LogoutOK, Meta: ctx.Payload.Meta}
}
