package handler

import (
	"encoding/json"
	"errors"
	infranats "gochat/internal/infra/nats"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

// BanUser - 群主將群內成員禁言
func BanUser(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.BanUserReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 || req.UserID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if errPayload := checkGroupOwner(ctx, req.GroupID, "BanUser"); errPayload != nil {
		return errPayload
	}

	if err := group.SetBan(req.GroupID, req.UserID, true); err != nil {
		switch {
		case errors.Is(err, group.ErrNotMember):
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		case errors.Is(err, group.ErrTargetOwner):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "cannot ban the owner", ctx.Payload)
		default:
			log.Printf("[BanUser] SetBan group=%d user=%d error: %v", req.GroupID, req.UserID, err)
			return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
		}
	}

	event := &protocol.BanUserEvent{GroupID: req.GroupID, UserID: req.UserID}
	if err := infranats.Publish(infranats.SubjectGroupBan, event); err != nil {
		log.Printf("[BanUser] Publish group=%d user=%d error: %v", req.GroupID, req.UserID, err)
	}

	log.Printf("[BanUser] OK owner=%d banned user=%d in group=%d", ctx.Client.UserId, req.UserID, req.GroupID)
	return &protocol.Payload{MsgType: protocol.BanUserOk, Meta: ctx.Payload.Meta}
}
