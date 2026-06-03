package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	infranats "gochat/internal/infra/nats"
	"log"

	"gochat/internal/chat"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
)

// requireGroupOwner 確認當前連線使用者正處於該群場景且為群主, 通過回傳 nil, 否則回傳對應錯誤 Payload。
func requireGroupOwner(ctx *ws.Ctx, groupID int, tag string) *protocol.Payload {
	sess := session.Get(ctx.Client.UserId, ctx.Client.ConnID)
	if sess.Scene != protocol.SceneInGroup || sess.InGroupID != groupID {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "not in group", ctx.Payload)
	}
	gu, err := group.GetMembership(sess.UserID, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		}
		log.Printf("[%s] GetMembership user=%d group=%d error: %v", tag, sess.UserID, groupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if gu.RoleType != group.RoleOwner {
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "not the owner", ctx.Payload)
	}
	return nil
}

func PinChat(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.PinChatReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 || req.ChatID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if errPayload := requireGroupOwner(ctx, req.GroupID, "PinChat"); errPayload != nil {
		return errPayload
	}

	ok, err := chat.ExistsInGroup(req.ChatID, req.GroupID)
	if err != nil {
		log.Printf("[PinChat] ExistsInGroup group=%d chat=%d error: %v", req.GroupID, req.ChatID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if !ok {
		return protocol.NewErrPayload(protocol.ErrorNotFound, "chat not found", ctx.Payload)
	}

	if _, err := group.SetPinChat(req.GroupID, req.ChatID); err != nil {
		log.Printf("[PinChat] SetPinChat group=%d chat=%d error: %v", req.GroupID, req.ChatID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	event := &protocol.PinChatCastEvent{GroupID: req.GroupID, ChatID: req.ChatID}
	if err := infranats.Publish(infranats.SubjectGroupPin, event); err != nil {
		log.Printf("[PinChat] Publish group=%d chat=%d error: %v", req.GroupID, req.ChatID, err)
	}

	return &protocol.Payload{MsgType: protocol.PinChatOk, Meta: ctx.Payload.Meta}
}
