package handler

import (
	"encoding/json"
	infranats "gochat/internal/infra/nats"
	"log"

	"gochat/internal/chat"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func PinChat(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.PinChatReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 || req.ChatID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if errPayload := checkGroupOwner(ctx, req.GroupID, "PinChat"); errPayload != nil {
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
