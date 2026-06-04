package handler

import (
	"encoding/json"
	infranats "gochat/internal/infra/nats"
	"log"

	"gochat/internal/chat"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func DelChat(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.DelChatReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 || req.ChatID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if errPayload := checkGroupOwner(ctx, req.GroupID, "DelChat"); errPayload != nil {
		return errPayload
	}

	n, err := chat.MarkDeleted(req.ChatID, req.GroupID)
	if err != nil {
		log.Printf("[DelChat] MarkDeleted group=%d chat=%d error: %v", req.GroupID, req.ChatID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if n == 0 {
		return protocol.NewErrPayload(protocol.ErrorNotFound, "chat not found", ctx.Payload)
	}

	event := &protocol.DelChatCastEvent{GroupID: req.GroupID, ChatID: req.ChatID}
	if err := infranats.Publish(infranats.SubjectGroupDel, event); err != nil {
		log.Printf("[DelChat] Publish group=%d chat=%d error: %v", req.GroupID, req.ChatID, err)
	}

	return &protocol.Payload{MsgType: protocol.DelChatOk, Meta: ctx.Payload.Meta}
}
