package handler

import (
	"encoding/json"
	"gochat/internal/protocol"
	"gochat/internal/user"
	"gochat/internal/ws"
)

func HandleLogin(ctx *ws.Ctx) *protocol.Payload {

	var req protocol.LoginReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.Token == "" {
		return protocol.NewErrPayload(protocol.Error, "invalid request", ctx.Payload)
	}

	tokenInfo, err := user.GetLoginToken(req)
	if err != nil {
		return protocol.NewErrPayload(protocol.Error, "invalid token", ctx.Payload)
	}

	u, err := user.Login(ctx, tokenInfo)
	if err != nil {
		return protocol.NewErrPayload(protocol.Error, "user login error", ctx.Payload)
	}

	if u.IsSuspended {
		return protocol.NewErrPayload(protocol.Error, "user is suspended", ctx.Payload)
	}

	ws.Register(ctx.Client)

	return &protocol.Payload{MsgType: protocol.LoginOk}
}
