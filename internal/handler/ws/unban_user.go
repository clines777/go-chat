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

// UnbanUser - 群主解除群內成員禁言
func UnbanUser(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.UnbanUserReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 || req.UserID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if errPayload := checkGroupOwner(ctx, req.GroupID, "UnbanUser"); errPayload != nil {
		return errPayload
	}

	if err := group.SetBan(req.GroupID, req.UserID, false); err != nil {
		switch {
		case errors.Is(err, group.ErrNotMember):
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		case errors.Is(err, group.ErrTargetOwner):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "cannot unban the owner", ctx.Payload)
		default:
			log.Printf("[UnbanUser] SetBan group=%d user=%d error: %v", req.GroupID, req.UserID, err)
			return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
		}
	}

	event := &protocol.UnbanUserEvent{GroupID: req.GroupID, UserID: req.UserID}
	if err := infranats.Publish(infranats.SubjectGroupUnban, event); err != nil {
		log.Printf("[UnbanUser] Publish group=%d user=%d error: %v", req.GroupID, req.UserID, err)
	}

	log.Printf("[UnbanUser] OK owner=%d unbanned user=%d in group=%d", ctx.Client.UserId, req.UserID, req.GroupID)
	return &protocol.Payload{MsgType: protocol.UnbanUserOk, Meta: ctx.Payload.Meta}
}
