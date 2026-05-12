package handler

import (
	"encoding/json"
	"gochat/internal/group"
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
	ws.Register(ctx.Client)

	userGroups, err := group.GetGroupsOfUser(u.ID, u.SiteBid)
	if err != nil {
		return protocol.NewErrPayload(protocol.Error, "fetch groups error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.LoginResp{
		UserID:     u.ID,
		Username:   u.ExtUsername,
		UserGroups: userGroups,
	})
	if err != nil {
		return protocol.NewErrPayload(protocol.Error, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.LoginOk, Data: respData, Meta: ctx.Payload.Meta}
}
