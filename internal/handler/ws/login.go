package handler

import (
	"encoding/json"
	"time"

	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func Login(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.LoginReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.Token == "" {
		return protocol.NewErrPayload(protocol.Error, "invalid token", ctx.Payload)
	}

	r := redis.GetRedis()
	key := protocol.LoginTokenKey(req.Token)

	var tokenData protocol.GetTokenReq
	ok, err := r.GetJSON(key, &tokenData)
	if err != nil || !ok {
		return protocol.NewErrPayload(protocol.Error, "token not found or expired", ctx.Payload)
	}

	_ = r.Del(key)

	sess := &ws.Session{
		ConnID:  ctx.Client.ConnID,
		UserID:  int64(tokenData.MemberId),
		SiteBid: tokenData.SiteBid,
	}
	if err := r.SetJSON(protocol.SessionKey(ctx.Client.ConnID), sess, 24*time.Hour); err != nil {
		return protocol.NewErrPayload(protocol.Error, "session error", ctx.Payload)
	}

	ws.Register(ctx.Client)

	return &protocol.Payload{MsgType: protocol.LoginOk}
}
