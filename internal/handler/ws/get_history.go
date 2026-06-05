package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"gochat/internal/chat"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

// GetHistory 向上分頁載入群組更舊的歷史訊息。
// 以 req.BeforeID (前端目前最舊一筆的 id) 為游標, 回傳 id 更小的前一頁訊息。
func GetHistory(ctx *ws.Ctx) *protocol.Payload {
	var req protocol.GetHistoryReq
	if err := json.Unmarshal(ctx.Payload.Data, &req); err != nil || req.GroupID == 0 {
		log.Printf("[GetHistory] bad request: %v", err)
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "invalid request", ctx.Payload)
	}

	sess := ctx.Session

	if _, err := group.GetMembership(sess.UserID, req.GroupID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		}
		log.Printf("[GetHistory] GetMembership user=%d group=%d error: %v", sess.UserID, req.GroupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	chats, hasMore, err := chat.GetHistoryChats(req.GroupID, req.BeforeID, RecentChatCount)
	if err != nil {
		log.Printf("[GetHistory] GetHistoryChats group=%d before=%d error: %v", req.GroupID, req.BeforeID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	respData, err := json.Marshal(&protocol.GetHistoryResp{
		GroupId: req.GroupID,
		Chats:   chats,
		HasMore: hasMore,
	})
	if err != nil {
		log.Printf("[GetHistory] marshal error: %v", err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}

	return &protocol.Payload{MsgType: protocol.GetHistoryOk, Data: respData, Meta: ctx.Payload.Meta}
}
