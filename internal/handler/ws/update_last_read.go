package handler

import (
	"encoding/json"
	"log"

	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func UpdateLastRead(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.UpdateLastReadReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID <= 0 || req.ChatID <= 0 {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	if err := group.UpdateLastRead(ctx.Client.UserId, req.GroupID, req.ChatID); err != nil {
		log.Printf("[UpdateLastRead] user=%d group=%d chat=%d error: %v", ctx.Client.UserId, req.GroupID, req.ChatID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.UpdateLastReadOk, Meta: ctx.Payload.Meta}
}
