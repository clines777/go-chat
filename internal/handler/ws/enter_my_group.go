package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

func EnterMyGroup(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.EnterMyGroupReq
	if ctx.Payload.Data != nil {
		if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil {
			log.Printf("[EnterMyGroup] bad request: %v", err)
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
		}
	}

	sess := session.Get(ctx.Client.ConnID)

	groups, err := group.GetMyGroupsPaged(sess.UserID, req.Page)
	if err != nil {
		log.Printf("[EnterMyGroup] GetMyGroupsPaged user=%d page=%d error: %v", sess.UserID, req.Page, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	ws.LeaveAllGroups(ctx.Client.ConnID)

	sess.Scene = protocol.SceneMyGroup
	sess.InGroupID = 0
	if err := session.Set(ctx.Client.ConnID, sess); err != nil {
		log.Printf("[EnterMyGroup] session.Set error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.EnterMyGroupResp{
		Page:   req.Page,
		Groups: groups,
	})
	if err != nil {
		log.Printf("[EnterMyGroup] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.EnterMyGroupOk, Data: respData, Meta: ctx.Payload.Meta}
}
