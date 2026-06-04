package handler

import (
	"encoding/json"
	infranats "gochat/internal/infra/nats"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

// UnpinChat 取消置頂消息
func UnpinChat(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.UnpinChatReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if errPayload := checkGroupOwner(ctx, req.GroupID, "UnpinChat"); errPayload != nil {
		return errPayload
	}

	if _, err := group.SetPinChat(req.GroupID, 0); err != nil {
		log.Printf("[UnpinChat] SetPinChat group=%d error: %v", req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	event := &protocol.UnpinChatCastEvent{GroupID: req.GroupID}
	if err := infranats.Publish(infranats.SubjectGroupUnpin, event); err != nil {
		log.Printf("[UnpinChat] Publish group=%d error: %v", req.GroupID, err)
	}

	return &protocol.Payload{MsgType: protocol.UnpinChatOk, Meta: ctx.Payload.Meta}
}
