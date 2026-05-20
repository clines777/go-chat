package handler

import (
	"encoding/json"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/user"
	"gochat/internal/ws"
)

func Login(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.LoginReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.Token == "" {
		return protocol.NewErrPayload(protocol.ErrWrongParam, "invalid request", ctx.Payload)
	}

	tokenInfo, err := user.GetLoginToken(req)
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "invalid token", ctx.Payload)
	}

	u, err := user.Login(ctx, tokenInfo)
	if err != nil {
		ctx.Client.Conn.Close()
		return protocol.NewErrPayload(protocol.ErrInternalError, "user login error", ctx.Payload)
	}
	ws.Register(ctx.Client)

	userGroups, err := group.GetMyGroups(u.ID, u.SiteBid)
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "fetch groups error", ctx.Payload)
	}

	apiToken, err := user.GenerateApiToken(u.ID, u.SiteBid)
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.LoginResp{
		UserID:     u.ID,
		Username:   u.ExtUsername,
		ApiToken:   apiToken,
		UserGroups: userGroups,
	})
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.LoginOK, Data: respData, Meta: ctx.Payload.Meta}
}
