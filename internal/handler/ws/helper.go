package handler

import (
	"database/sql"
	"errors"
	"gochat/internal/chat"
	"gochat/internal/group"
	"gochat/internal/protocol"
	"gochat/internal/ws"
	"log"
)

// checkGroupOwner 確認當前連線使用者正處於該群場景且為群主, 通過回傳 nil, 否則回傳對應錯誤 Payload。
func checkGroupOwner(ctx *ws.Ctx, groupID int, tag string) *protocol.Payload {
	sess := ctx.Session
	if sess.Scene != protocol.SceneInGroup || sess.InGroupID != groupID {
		return protocol.NewErrPayload(protocol.ErrInvalidParam, "not in group", ctx.Payload)
	}
	gu, err := group.GetMembership(sess.UserID, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.NewErrPayload(protocol.ErrNotMember, "not a member", ctx.Payload)
		}
		log.Printf("[%s] getting membership: user=%d group=%d error: %v", tag, sess.UserID, groupID, err)
		return protocol.NewErrPayload(protocol.ErrInternalError, "internal error", ctx.Payload)
	}
	if gu.RoleType != group.RoleOwner {
		return protocol.NewErrPayload(protocol.ErrUnauthorized, "not the owner", ctx.Payload)
	}
	return nil
}

// FindPinChat 解析群組置頂訊息: 回傳有效的 pinChatID (置頂訊息已刪除或不存在時為 0) 與訊息本體。
func FindPinChat(groupID, pinChatID int) *protocol.ChatInfo {
	if pinChatID <= 0 {
		return nil
	}
	c, err := chat.FindInGroup(pinChatID, groupID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("FindPinChat group=%d pin=%d error: %v", groupID, pinChatID, err)
		}
		return nil
	}
	return c
}
