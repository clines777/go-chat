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

// KickUser - 群主將成員移出群組
func KickUser(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.KickUserReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 || req.UserID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if errPayload := checkGroupOwner(ctx, req.GroupID, "KickUser"); errPayload != nil {
		return errPayload
	}

	if err := group.Kick(req.GroupID, req.UserID); err != nil {
		switch {
		case errors.Is(err, group.ErrNotMember):
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		case errors.Is(err, group.ErrTargetOwner):
			return protocol.NewErrPayload(protocol.ErrInvalidParam, "cannot kick the owner", ctx.Payload)
		default:
			log.Printf("[KickUser] Kick group=%d user=%d error: %v", req.GroupID, req.UserID, err)
			return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
		}
	}

	event := &protocol.KickUserEvent{GroupID: req.GroupID, UserID: req.UserID}
	if err := infranats.Publish(infranats.SubjectGroupKick, event); err != nil {
		log.Printf("[KickUser] Publish group=%d user=%d error: %v", req.GroupID, req.UserID, err)
	}

	log.Printf("[KickUser] OK owner=%d kicked user=%d from group=%d", ctx.Client.UserId, req.UserID, req.GroupID)
	return &protocol.Payload{MsgType: protocol.KickUserOk, Meta: ctx.Payload.Meta}
}
