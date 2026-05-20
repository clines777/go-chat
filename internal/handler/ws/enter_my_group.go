package handler

import (
	"encoding/json"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

func EnterMyGroup(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.EnterMyGroupReq
	if ctx.Payload.Data != nil {
		if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil {
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
		}
	}

	sess := session.Get(ctx.Client.ConnID)

	groups, err := group.GetMyGroupsPaged(sess.UserID, sess.SiteBid, req.Page)
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	sess.Scene = protocol.SceneMyGroup
	sess.InGroupId = 0
	if err := session.Set(ctx.Client.ConnID, sess); err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.EnterMyGroupResp{
		Page:   req.Page,
		Groups: groups,
	})
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.EnterMyGroupOk, Data: respData, Meta: ctx.Payload.Meta}
}
