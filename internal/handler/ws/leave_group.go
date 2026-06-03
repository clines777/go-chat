package handler

import (
	"encoding/json"
	"errors"
	infranats "gochat/internal/infra/nats"
	"log"
	"time"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

// LeaveGroup - 用戶退出所屬群,
func LeaveGroup(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.LeaveGroupReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID <= 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)

	if err := group.Leave(sess.UserID, req.GroupID); err != nil {
		switch {
		case errors.Is(err, group.ErrNotMember):
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		case errors.Is(err, group.ErrIsOwner):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "群主無法離開群組，請先解散或轉讓群組", ctx.Payload)
		default:
			log.Printf("[LeaveGroup] Leave user=%d group=%d error: %v", sess.UserID, req.GroupID, err)
			return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
		}
	}

	ws.ResetGroupScene(ctx.Client.ConnID)

	sess.InGroupID = 0
	sess.Scene = protocol.SceneMyGroup
	if err := session.Set(ctx.Client.UserId, ctx.Client.ConnID, sess); err != nil {
		log.Printf("[LeaveGroup] session.Set error: %v", err)
	}

	respData, err := json.Marshal(map[string]int{"group_id": req.GroupID})
	if err != nil {
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	event := &protocol.LeaveGroupEvent{GroupId: req.GroupID, UserId: sess.UserID, Time: time.Now().Unix()}
	err = infranats.Publish(infranats.SubjectGroupLeave, event)
	if err != nil {
		log.Printf("[LeaveGroup] Publish error: %v", err)
	}

	log.Printf("[LeaveGroup] OK user=%d left group=%d", sess.UserID, req.GroupID)
	return &protocol.Payload{MsgType: protocol.LeaveGroupOk, Data: respData, Meta: ctx.Payload.Meta}
}
