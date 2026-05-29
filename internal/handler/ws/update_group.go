package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	infranats "gochat/internal/infra/nats"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

func UpdateGroup(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.UpdateGroupReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if req.Title == "" && req.Bulletin == "" && req.Remark == "" {
		return &protocol.Payload{MsgType: protocol.UpdateGroupOk, Meta: ctx.Payload.Meta}
	}

	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)

	gu, err := group.GetMembership(sess.UserID, req.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		}
		log.Printf("[UpdateGroup] GetMembership user=%d group=%d error: %v", sess.UserID, req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if gu.RoleType != group.RoleOwner {
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "not the owner", ctx.Payload)
	}

	g, err := group.FindByID(req.GroupID)
	if err != nil {
		log.Printf("[UpdateGroup] FindByID group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	updateCount := 0
	if req.Title != "" && req.Title != g.Title || req.Bulletin != g.Bulletin || req.Remark != g.Remark {
		updateCount, err = group.Update(req.GroupID, req.Title, req.Bulletin, req.Remark)
		if err != nil {
			log.Printf("[UpdateGroup] Update group=%d error: %v", req.GroupID, err)
			return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
		}
	}

	if updateCount > 0 {
		event := &protocol.UpdateGroupCastEvent{GroupID: req.GroupID, Title: req.Title, Bulletin: req.Bulletin, Remark: req.Remark}
		if err = infranats.Publish(fmt.Sprintf("group.update.%d", req.GroupID), event); err != nil {
			log.Printf("[UpdateGroup] Publish group=%d error: %v", req.GroupID, err)
		}
	}

	return &protocol.Payload{MsgType: protocol.UpdateGroupOk, Meta: ctx.Payload.Meta}
}
