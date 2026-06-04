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

// DelChat 群主刪除群內留言
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

	// 若刪到的正好是置頂發言, 順手取消置頂並通知群內
	if cleared, err := group.ClearPinChatIfMatch(req.GroupID, req.ChatID); err != nil {
		log.Printf("[DelChat] ClearPinChatIfMatch group=%d chat=%d error: %v", req.GroupID, req.ChatID, err)
	} else if cleared > 0 {
		unpinEvent := &protocol.UnpinChatCastEvent{GroupID: req.GroupID}
		if err := infranats.Publish(infranats.SubjectGroupUnpin, unpinEvent); err != nil {
			log.Printf("[DelChat] Publish unpin group=%d error: %v", req.GroupID, err)
		}
	}

	return &protocol.Payload{MsgType: protocol.DelChatOk, Meta: ctx.Payload.Meta}
}
